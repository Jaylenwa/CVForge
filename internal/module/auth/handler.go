package auth

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"openresume/internal/infra/config"
	conf "openresume/internal/module/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	svc *Service
}

func NewHandler(cfg config.Config, sysConfig *conf.Service, rdb *redis.Client, db *gorm.DB) *Handler {
	return &Handler{svc: NewService(cfg, sysConfig, rdb, db)}
}

func popupHTML(origin, access, refresh string, user map[string]any) string {
	b, _ := json.Marshal(map[string]any{"status": "ok", "accessToken": access, "refreshToken": refresh, "user": user})
	return "<!doctype html><html><head><meta charset=\"utf-8\"><title>WeChat Login</title></head><body><script>(function(){try{var data=" + string(b) + ";window.opener&&window.opener.postMessage(data,'" + origin + "');window.close();}catch(e){document.body.innerText='Login succeeded, but messaging failed';}})();</script></body></html>"
}

func (h *Handler) WeChatRedirect(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.FeatureWeChatLogin == "off" {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "feature disabled"})
			return
		}
		client := c.Query("client")
		if client != "popup" && client != "redirect" {
			client = "popup"
		}
		origin := c.Query("origin")
		if origin != "" && !h.svc.IsAllowedOrigin(cfg.OAuthAllowedOrigins, origin) {
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
		_ = h.svc.SaveOAuthState(state, st)
		redir := h.svc.MakeWeChatLoginURL(state)
		if redir == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "wechat not configured"})
			return
		}
		c.Redirect(http.StatusFound, redir)
	}
}

func (h *Handler) GithubRedirect(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.FeatureGithubLogin == "off" {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "feature disabled"})
			return
		}
		client := c.Query("client")
		if client != "popup" && client != "redirect" {
			client = "popup"
		}
		origin := c.Query("origin")
		if origin != "" && !h.svc.IsAllowedOrigin(cfg.OAuthAllowedOrigins, origin) {
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
		_ = h.svc.SaveOAuthState(state, st)
		redir := h.svc.MakeGithubLoginURL(state)
		if redir == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "github not configured"})
			return
		}
		c.Redirect(http.StatusFound, redir)
	}
}

func (h *Handler) GithubCallback(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.FeatureGithubLogin == "off" {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "feature disabled"})
			return
		}
		code := c.Query("code")
		state := c.Query("state")
		if code == "" || state == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing code/state"})
			return
		}
		val, err := h.svc.GetOAuthState(state)
		if err != nil || val == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
			return
		}
		h.svc.DelOAuthState(state)
		var st struct {
			Client string `json:"client"`
			Origin string `json:"origin"`
			IP     string `json:"ip"`
			UA     string `json:"ua"`
		}
		_ = json.Unmarshal([]byte(val), &st)
		tokenResp, err := h.svc.ExchangeGithubCode(code)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "exchange failed"})
			return
		}
		profile, _ := h.svc.FetchGithubUserInfo(tokenResp.AccessToken)
		user, err := h.svc.FindOrCreateGithubUser(profile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "account error"})
			return
		}
		access, refresh := h.svc.IssueTokens(user.ID)
		if st.Client == "popup" {
			origin := st.Origin
			if origin == "" {
				origin = cfg.FrontendBaseURL
			}
			html := popupHTML(origin, access, refresh, h.svc.SanitizeUser(user))
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
			return
		}
		ott := uuid.NewString()
		payload := gin.H{"accessToken": access, "refreshToken": refresh, "user": h.svc.SanitizeUser(user)}
		_ = h.svc.SaveOTT(ott, payload)
		u := cfg.FrontendBaseURL
		if u == "" {
			u = "http://localhost:3000"
		}
		c.Redirect(http.StatusFound, u+"/#/oauth/callback?ott="+url.QueryEscape(ott))
	}
}

func (h *Handler) WeChatCallback(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		val, err := h.svc.GetOAuthState(state)
		if err != nil || val == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
			return
		}
		h.svc.DelOAuthState(state)
		var st struct {
			Client string `json:"client"`
			Origin string `json:"origin"`
			IP     string `json:"ip"`
			UA     string `json:"ua"`
		}
		_ = json.Unmarshal([]byte(val), &st)
		tokenResp, err := h.svc.ExchangeCode(code)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "exchange failed"})
			return
		}
		profile, _ := h.svc.FetchUserInfo(tokenResp.AccessToken, tokenResp.OpenID)
		user, err := h.svc.FindOrCreateWeChatUser(profile, tokenResp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "account error"})
			return
		}
		access, refresh := h.svc.IssueTokens(user.ID)
		if st.Client == "popup" {
			origin := st.Origin
			if origin == "" {
				origin = cfg.FrontendBaseURL
			}
			html := popupHTML(origin, access, refresh, h.svc.SanitizeUser(user))
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
			return
		}
		ott := uuid.NewString()
		payload := gin.H{"accessToken": access, "refreshToken": refresh, "user": h.svc.SanitizeUser(user)}
		_ = h.svc.SaveOTT(ott, payload)
		u := cfg.FrontendBaseURL
		if u == "" {
			u = "http://localhost:3000"
		}
		c.Redirect(http.StatusFound, u+"/#/oauth/callback?ott="+url.QueryEscape(ott))
	}
}

func (h *Handler) ConsumeOTT(c *gin.Context) {
	var body struct {
		OTT string `json:"ott"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.OTT == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	val, err := h.svc.GetOTT(body.OTT)
	if err != nil || val == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ott"})
		return
	}
	h.svc.DelOTT(body.OTT)
	var payload map[string]any
	_ = json.Unmarshal([]byte(val), &payload)
	c.JSON(http.StatusOK, payload)
}

func (h *Handler) SendCode(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}
	code := h.svc.GenerateVerifyCode()
	_ = h.svc.SaveVerifyCode(req.Email, code)
	if err := h.svc.SendCode(req.Email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "send failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) Register(c *gin.Context) {
	var req struct{ Email, Code, Password, Name string }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	access, refresh, err := h.svc.Register(req.Email, req.Code, req.Password, req.Name)
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			c.JSON(http.StatusConflict, gin.H{"error": "email exists"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accessToken": access, "refreshToken": refresh})
}

func (h *Handler) Login(c *gin.Context) {
	var req struct{ Email, Password string }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	access, refresh, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accessToken": access, "refreshToken": refresh})
}

func (h *Handler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	access, refresh, err := h.svc.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accessToken": access, "refreshToken": refresh})
}

func (h *Handler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := h.svc.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
