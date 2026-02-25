package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"openresume/internal/infra/config"
	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		if t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("invalid signing algorithm")
			}
			return []byte(config.CF.JWTSecret), nil
		}); err == nil && t.Valid {
			if claims, ok := t.Claims.(jwt.MapClaims); ok {
				if uidF, ok2 := claims["uid"].(float64); ok2 && uidF > 0 {
					_ = h.svc.GlobalLogout(uint(uidF))
				}
			}
		}
		_ = h.svc.LogoutAccess(strings.TrimPrefix(auth, "Bearer "))
	}
	if err := h.svc.Logout(req.RefreshToken); err != nil {
		logger.WithCtx(c).Error("auth.logout failed", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
