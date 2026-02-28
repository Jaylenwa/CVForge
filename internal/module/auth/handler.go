package auth

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"
	"cvforge/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	svc *Service
}

func NewHandler() *Handler {
	return &Handler{svc: NewService()}
}

func popupHTML(origin, access, refresh string, user map[string]any) string {
	b, _ := json.Marshal(map[string]any{"status": "ok", "accessToken": access, "refreshToken": refresh, "user": user})
	return "<!doctype html><html><head><meta charset=\"utf-8\"><title>WeChat Login</title></head><body><script>(function(){try{var data=" + string(b) + ";window.opener&&window.opener.postMessage(data,'" + origin + "');window.close();}catch(e){document.body.innerText='Login succeeded, but messaging failed';}})();</script></body></html>"
}

func (h *Handler) GithubRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !h.svc.FeatureGithubEnabled() {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "feature disabled"})
			return
		}
		client := c.Query("client")
		if client != "popup" && client != "redirect" {
			client = "popup"
		}
		origin := c.Query("origin")
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

func (h *Handler) GithubCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !h.svc.FeatureGithubEnabled() {
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
			logger.WithCtx(c).Error("auth.github invalid state", zap.Error(err))
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
			logger.WithCtx(c).Error("auth.github exchange failed", zap.Error(err))
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		profile, _ := h.svc.FetchGithubUserInfo(tokenResp.AccessToken)
		user, err := h.svc.FindOrCreateGithubUser(profile)
		if err != nil {
			logger.WithCtx(c).Error("auth.github account error", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.svc.TouchLastLoginAt(user.ID)
		access, refresh := h.svc.IssueTokens(user.ID)
		if st.Client == "popup" {
			origin := st.Origin
			if origin == "" {
				origin = h.svc.FrontendBase()
			}
			html := popupHTML(origin, access, refresh, h.svc.SanitizeUser(user))
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
			return
		}
		ott := uuid.NewString()
		payload := gin.H{"accessToken": access, "refreshToken": refresh, "user": h.svc.SanitizeUser(user)}
		_ = h.svc.SaveOTT(ott, payload)
		u := h.svc.FrontendBase()
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
		logger.WithCtx(c).Error("auth.consume_ott bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	val, err := h.svc.GetOTT(body.OTT)
	if err != nil || val == "" {
		logger.WithCtx(c).Error("auth.consume_ott invalid ott", zap.Error(err))
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
		logger.WithCtx(c).Error("auth.send_code bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}
	code := h.svc.GenerateVerifyCode()
	_ = h.svc.SaveVerifyCode(req.Email, code)
	if err := h.svc.SendCode(req.Email, code); err != nil {
		logger.WithCtx(c).Error("auth.send_code failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) Register(c *gin.Context) {
	var req struct{ Email, Code, Password, Name string }
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithCtx(c).Error("auth.register bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	access, refresh, err := h.svc.Register(req.Email, req.Code, req.Password, req.Name)
	if err != nil {
		logger.WithCtx(c).Error("auth.register failed", zap.Error(err))
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
		logger.WithCtx(c).Error("auth.login bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	access, refresh, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		logger.WithCtx(c).Error("auth.login failed", zap.Error(err))
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
		logger.WithCtx(c).Error("auth.refresh bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	access, refresh, err := h.svc.Refresh(req.RefreshToken)
	if err != nil {
		logger.WithCtx(c).Error("auth.refresh failed", zap.Error(err))
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
		logger.WithCtx(c).Error("auth.logout bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	authHeader := c.GetHeader("Authorization")
	accessToken := ""
	if strings.HasPrefix(authHeader, "Bearer ") {
		accessToken = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	}
	if err := h.svc.LogoutSession(accessToken, req.RefreshToken); err != nil {
		logger.WithCtx(c).Error("auth.logout failed", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

type weChatMPEvent struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Event        string   `xml:"Event"`
	EventKey     string   `xml:"EventKey"`
	Ticket       string   `xml:"Ticket"`
}

type weChatMPEncryptWrap struct {
	XMLName xml.Name `xml:"xml"`
	Encrypt string   `xml:"Encrypt"`
}

func (h *Handler) WeChatMPCreateScene(c *gin.Context) {
	if !h.svc.FeatureWeChatMPEnabled() {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "feature disabled"})
		return
	}
	scene, qrURL, expireSeconds, err := h.svc.CreateWeChatMPLoginScene(context.Background())
	if err != nil {
		logger.WithCtx(c).Error("auth.wechatmp create scene failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"scene": scene, "qrUrl": qrURL, "expiresIn": expireSeconds})
}

func (h *Handler) WeChatMPSceneStatus(c *gin.Context) {
	if !h.svc.FeatureWeChatMPEnabled() {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "feature disabled"})
		return
	}
	scene := c.Param("scene")
	payload, err := h.svc.GetWeChatMPScenePayload(context.Background(), scene)
	if err != nil || payload.Scene == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": "expired"})
		return
	}
	if payload.Status == "ok" && payload.OTT != "" {
		if cache.RDB != nil {
			_ = cache.RDB.Del(context.Background(), common.RedisKeyWeChatMPScene.F(scene)).Err()
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "ott": payload.OTT})
		return
	}
	if payload.ExpireAt > 0 && time.Now().Unix() >= payload.ExpireAt {
		if cache.RDB != nil {
			_ = cache.RDB.Del(context.Background(), common.RedisKeyWeChatMPScene.F(scene)).Err()
		}
		c.JSON(http.StatusNotFound, gin.H{"status": "expired"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "pending"})
}

func (h *Handler) WeChatMPCallbackGet(c *gin.Context) {
	if !h.svc.FeatureWeChatMPEnabled() {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "feature disabled"})
		return
	}
	msgSignature := c.Query("msg_signature")
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")
	if echostr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		return
	}
	if msgSignature != "" {
		if !h.svc.VerifyWeChatMPMsgSignature(msgSignature, timestamp, nonce, echostr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			return
		}
		plain, err := h.svc.decryptWeChatMPText(echostr)
		if err != nil || plain == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			return
		}
		c.String(http.StatusOK, plain)
		return
	}
	if !h.svc.VerifyWeChatMPSignature(signature, timestamp, nonce) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		return
	}
	c.String(http.StatusOK, echostr)
}

func (h *Handler) WeChatMPCallbackPost(c *gin.Context) {
	if !h.svc.FeatureWeChatMPEnabled() {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "feature disabled"})
		return
	}
	msgSignature := c.Query("msg_signature")
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")

	raw, err := c.GetRawData()
	if err != nil || len(raw) == 0 {
		c.String(http.StatusOK, "success")
		return
	}
	if msgSignature != "" {
		var wrap weChatMPEncryptWrap
		if e := xml.Unmarshal(raw, &wrap); e != nil || strings.TrimSpace(wrap.Encrypt) == "" {
			c.String(http.StatusOK, "success")
			return
		}
		enc := strings.TrimSpace(wrap.Encrypt)
		if !h.svc.VerifyWeChatMPMsgSignature(msgSignature, timestamp, nonce, enc) {
			c.String(http.StatusOK, "success")
			return
		}
		plain, e := h.svc.decryptWeChatMPText(enc)
		if e != nil || plain == "" {
			c.String(http.StatusOK, "success")
			return
		}
		raw = []byte(plain)
	} else {
		if !h.svc.VerifyWeChatMPSignature(signature, timestamp, nonce) {
			c.String(http.StatusOK, "success")
			return
		}
	}

	var ev weChatMPEvent
	if e := xml.Unmarshal(raw, &ev); e != nil {
		c.String(http.StatusOK, "success")
		return
	}

	openid := strings.TrimSpace(ev.FromUserName)
	scene := ""
	switch strings.ToUpper(strings.TrimSpace(ev.Event)) {
	case "SUBSCRIBE":
		if strings.HasPrefix(ev.EventKey, "qrscene_") {
			scene = strings.TrimPrefix(ev.EventKey, "qrscene_")
		}
	case "SCAN":
		scene = strings.TrimSpace(ev.EventKey)
	}
	if openid == "" || scene == "" {
		c.String(http.StatusOK, "success")
		return
	}

	p, e := h.svc.GetWeChatMPScenePayload(context.Background(), scene)
	if e != nil || p.Scene == "" || p.Status != "pending" {
		c.String(http.StatusOK, "success")
		return
	}

	ui, _ := h.svc.FetchWeChatMPUserInfo(context.Background(), openid)
	user, err := h.svc.FindOrCreateWeChatMPUser(ui, openid)
	if err != nil {
		logger.WithCtx(c).Error("auth.wechatmp account error", zap.Error(err))
		c.String(http.StatusOK, "success")
		return
	}
	if _, e2 := h.svc.MarkWeChatMPSceneReady(context.Background(), scene, openid, user); e2 != nil {
		logger.WithCtx(c).Error("auth.wechatmp mark scene ready failed", zap.Error(e2))
	}
	c.String(http.StatusOK, "success")
}
