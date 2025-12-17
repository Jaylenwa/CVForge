package handlers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"openresume/config"
	"openresume/mailer"
	"openresume/middleware"
	"openresume/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type sendCodeReq struct {
	Email string `json:"email"`
}
type registerReq struct{ Email, Code, Password, Name string }
type loginReq struct{ Email, Password string }
type refreshReq struct {
	RefreshToken string `json:"refreshToken"`
}

func RegisterAuthRoutes(r *gin.RouterGroup, cfg config.Config, rdb *redis.Client, db *gorm.DB) {
	auth := r.Group("/auth")
	// WeChat OAuth redirect
	auth.GET("/wechat/redirect", middleware.RateLimit(rdb, 10, time.Minute), func(c *gin.Context) {
		if cfg.FeatureWeChatLogin == "off" {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "feature disabled"})
			return
		}
		client := c.Query("client")
		if client != "popup" && client != "redirect" {
			client = "popup"
		}
		origin := c.Query("origin")
		if origin != "" && !isAllowedOrigin(cfg.OAuthAllowedOrigins, origin) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "origin not allowed"})
			return
		}
		state := uuid.NewString()
		st := map[string]any{
			"client":  client,
			"origin":  origin,
			"ip":      c.ClientIP(),
			"ua":      c.GetHeader("User-Agent"),
			"created": time.Now().Unix(),
		}
		b, _ := json.Marshal(st)
		_ = rdb.Set(context.Background(), "oauth:state:"+state, string(b), 10*time.Minute).Err()
		redir := makeWeChatLoginURL(cfg, state)
		if redir == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "wechat not configured"})
			return
		}
		c.Redirect(http.StatusFound, redir)
	})
	// WeChat OAuth callback
	auth.GET("/wechat/callback", middleware.RateLimit(rdb, 30, time.Minute), func(c *gin.Context) {
		if cfg.FeatureWeChatLogin == "off" {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "feature disabled"})
			return
		}
		code := c.Query("code")
		state := c.Query("state")
		if code == "" || state == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing code/state"})
			return
		}
		val, err := rdb.Get(context.Background(), "oauth:state:"+state).Result()
		if err != nil || val == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
			return
		}
		_ = rdb.Del(context.Background(), "oauth:state:"+state).Err()
		var st struct {
			Client string `json:"client"`
			Origin string `json:"origin"`
			IP     string `json:"ip"`
			UA     string `json:"ua"`
		}
		_ = json.Unmarshal([]byte(val), &st)
		tokenResp, err := exchangeCode(cfg, code)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "exchange failed"})
			return
		}
		profile, _ := fetchUserInfo(tokenResp.AccessToken, tokenResp.OpenID)
		user, err := findOrCreateWeChatUser(db, profile, tokenResp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "account error"})
			return
		}
		access, refresh := issueTokens(cfg, user.ID)
		if st.Client == "popup" {
			origin := st.Origin
			if origin == "" {
				origin = cfg.FrontendBaseURL
			}
			html := popupHTML(origin, access, refresh, user)
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
			return
		}
		ott := uuid.NewString()
		payload := gin.H{"accessToken": access, "refreshToken": refresh, "user": sanitizeUser(user)}
		pb, _ := json.Marshal(payload)
		_ = rdb.Set(context.Background(), "oauth:ott:"+ott, string(pb), time.Minute).Err()
		u := cfg.FrontendBaseURL
		if u == "" {
			u = "http://localhost:3000"
		}
		c.Redirect(http.StatusFound, u+"/#/oauth/callback?ott="+url.QueryEscape(ott))
	})
	// Consume one-time token to retrieve JWTs
	auth.POST("/wechat/consume-ott", func(c *gin.Context) {
		var body struct {
			OTT string `json:"ott"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.OTT == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		val, err := rdb.Get(context.Background(), "oauth:ott:"+body.OTT).Result()
		if err != nil || val == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ott"})
			return
		}
		_ = rdb.Del(context.Background(), "oauth:ott:"+body.OTT).Err()
		var payload map[string]any
		_ = json.Unmarshal([]byte(val), &payload)
		c.JSON(http.StatusOK, payload)
	})
	auth.POST("/send-code", middleware.RateLimit(rdb, 3, time.Minute), func(c *gin.Context) {
		var req sendCodeReq
		if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
			return
		}
		code := func() string {
			n := 6
			out := make([]byte, n)
			for i := 0; i < n; i++ {
				v, _ := rand.Int(rand.Reader, big.NewInt(10))
				out[i] = byte('0' + v.Int64())
			}
			return string(out)
		}()
		_ = rdb.Set(context.Background(), "verify:"+req.Email, code, 10*time.Minute).Err()
		if err := mailer.SendVerificationCode(cfg, req.Email, code); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "send failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	auth.POST("/register", middleware.RateLimit(rdb, 5, time.Minute), func(c *gin.Context) {
		var req registerReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		val, err := rdb.Get(context.Background(), "verify:"+req.Email).Result()
		if err != nil || val != req.Code {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
			return
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		u := models.User{Email: req.Email, PasswordHash: string(hash), Name: req.Name}
		if err := db.Create(&u).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email exists"})
			return
		}
		access, refresh := issueTokens(cfg, u.ID)
		c.JSON(http.StatusOK, gin.H{"accessToken": access, "refreshToken": refresh})
	})

	auth.POST("/login", middleware.RateLimit(rdb, 5, time.Minute), func(c *gin.Context) {
		var req loginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		var u models.User
		if err := db.Where("email = ?", req.Email).First(&u).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		access, refresh := issueTokens(cfg, u.ID)
		c.JSON(http.StatusOK, gin.H{"accessToken": access, "refreshToken": refresh})
	})

	auth.POST("/refresh", func(c *gin.Context) {
		var req refreshReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		t, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) { return []byte(cfg.JWTSecret), nil })
		if err != nil || !t.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, _ := t.Claims.(jwt.MapClaims)
		jti, _ := claims["jti"].(string)
		if jti != "" {
			if rdb.Get(context.Background(), "jwt:blacklist:"+jti).Val() == "1" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "revoked"})
				return
			}
		}
		uid := uint(claims["uid"].(float64))
		access, refresh := issueTokens(cfg, uid)
		c.JSON(http.StatusOK, gin.H{"accessToken": access, "refreshToken": refresh})
	})

	auth.POST("/logout", func(c *gin.Context) {
		var req refreshReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		t, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) { return []byte(cfg.JWTSecret), nil })
		if err != nil || !t.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, _ := t.Claims.(jwt.MapClaims)
		jti, _ := claims["jti"].(string)
		if jti != "" {
			_ = rdb.Set(context.Background(), "jwt:blacklist:"+jti, "1", time.Hour*24*7).Err()
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
}

func issueTokens(cfg config.Config, uid uint) (string, string) {
	mk := func(exp time.Duration) string {
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"uid": uid, "exp": time.Now().Add(exp).Unix(), "jti": uuid.NewString()})
		s, _ := t.SignedString([]byte(cfg.JWTSecret))
		return s
	}
	return mk(2 * time.Hour), mk(7 * 24 * time.Hour)
}

type wechatTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}
type wechatUserInfo struct {
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	HeadImgURL string `json:"headimgurl"`
	UnionID    string `json:"unionid"`
}

func makeWeChatLoginURL(cfg config.Config, state string) string {
	if cfg.WeChatAppID == "" || cfg.WeChatRedirectURI == "" {
		return ""
	}
	params := url.Values{}
	params.Set("appid", cfg.WeChatAppID)
	params.Set("redirect_uri", cfg.WeChatRedirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "snsapi_login")
	params.Set("state", state)
	return "https://open.weixin.qq.com/connect/qrconnect?" + params.Encode() + "#wechat_redirect"
}

func weChatExchangeCode(cfg config.Config, code string) (wechatTokenResponse, error) {
	var out wechatTokenResponse
	if cfg.WeChatAppID == "" || cfg.WeChatAppSecret == "" {
		return out, http.ErrNotSupported
	}
	u := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + url.QueryEscape(cfg.WeChatAppID) +
		"&secret=" + url.QueryEscape(cfg.WeChatAppSecret) +
		"&code=" + url.QueryEscape(code) +
		"&grant_type=authorization_code"
	resp, err := http.Get(u)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out.ErrCode != 0 || out.OpenID == "" || out.AccessToken == "" {
		return out, http.ErrHandlerTimeout
	}
	return out, nil
}

func weChatFetchUserInfo(accessToken, openid string) (wechatUserInfo, error) {
	var out wechatUserInfo
	if accessToken == "" || openid == "" {
		return out, http.ErrNotSupported
	}
	u := "https://api.weixin.qq.com/sns/userinfo?access_token=" + url.QueryEscape(accessToken) + "&openid=" + url.QueryEscape(openid)
	resp, err := http.Get(u)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out.OpenID == "" {
		return out, http.ErrHandlerTimeout
	}
	return out, nil
}

var exchangeCode = weChatExchangeCode
var fetchUserInfo = weChatFetchUserInfo

func popupHTML(origin, access, refresh string, user models.User) string {
	payload := map[string]any{"status": "ok", "accessToken": access, "refreshToken": refresh, "user": sanitizeUser(user)}
	b, _ := json.Marshal(payload)
	return "<!doctype html><html><head><meta charset=\"utf-8\"><title>WeChat Login</title></head><body><script>(function(){try{var data=" + string(b) + ";window.opener&&window.opener.postMessage(data,'" + origin + "');window.close();}catch(e){document.body.innerText='Login succeeded, but messaging failed';}})();</script></body></html>"
}

func findOrCreateWeChatUser(db *gorm.DB, ui wechatUserInfo, tr wechatTokenResponse) (models.User, error) {
	var user models.User
	var oa models.OAuthAccount
	q := db.Where("provider = ? AND provider_union_id <> '' AND provider_union_id = ?", "wechat", ui.UnionID).First(&oa)
	if ui.UnionID != "" && q.Error == nil {
		if err := db.First(&user, oa.UserID).Error; err == nil {
			return user, nil
		}
	}
	if err := db.Where("provider = ? AND provider_open_id = ?", "wechat", tr.OpenID).First(&oa).Error; err == nil {
		if err := db.First(&user, oa.UserID).Error; err == nil {
			return user, nil
		}
	}
	user = models.User{
		Name:      ui.Nickname,
		AvatarURL: ui.HeadImgURL,
		IsActive:  true,
		Role:      "user",
	}
	if err := db.Create(&user).Error; err != nil {
		return models.User{}, err
	}
	record := models.OAuthAccount{
		UserID:          user.ID,
		Provider:        "wechat",
		ProviderOpenID:  tr.OpenID,
		ProviderUnionID: ui.UnionID,
	}
	if b, e := json.Marshal(ui); e == nil {
		record.RawProfileJSON = string(b)
	}
	_ = db.Create(&record).Error
	return user, nil
}

func sanitizeUser(u models.User) map[string]any {
	return map[string]any{
		"id":        u.ID,
		"email":     u.Email,
		"name":      u.Name,
		"avatarUrl": u.AvatarURL,
		"language":  u.Language,
		"role":      u.Role,
	}
}

func isAllowedOrigin(origins string, origin string) bool {
	if origin == "" {
		return false
	}
	for _, o := range splitCSV(origins) {
		if o == origin {
			return true
		}
	}
	return false
}

func splitCSV(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if start < i {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	return out
}
