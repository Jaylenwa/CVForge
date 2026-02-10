package share

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"openresume/internal/common"
	"openresume/internal/infra/cache"
	"openresume/internal/middleware"
	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	svc *Service
}

func NewHandler() *Handler {
	return &Handler{svc: NewService()}
}

func (h *Handler) PublishResume(c *gin.Context) {
	uid, ok := middleware.UID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	sl, code, err := h.svc.PublishResumeForUser(uid, uint(id))
	if err != nil {
		logger.WithCtx(c).Error("share.publish failed", zap.Error(err), zap.Int("code", code), zap.String("id", c.Param("id")))
		switch code {
		case 403:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "resume not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"slug":    sl.Slug,
		"url":     "/public/resumes/" + sl.Slug,
		"apiUrl":  "/public/resumes/" + sl.Slug,
		"pageUrl": "/#/public/" + sl.Slug,
	})
}

func (h *Handler) GetPublic(c *gin.Context) {
	slug := c.Param("slug")
	shareToken := strings.TrimSpace(c.GetHeader("X-Share-Token"))
	if shareToken == "" {
		shareToken = strings.TrimSpace(c.Query("token"))
	}
	val, code, err := h.svc.GetPublicPayload(slug, shareToken)
	if err != nil {
		logger.WithCtx(c).Error("share.get_public failed", zap.Error(err), zap.Int("code", code), zap.String("slug", slug))
		switch code {
		case 401:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "password_required"})
		case 410:
			c.JSON(http.StatusGone, gin.H{"error": "expired"})
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		}
		return
	}
	c.Data(http.StatusOK, "application/json", []byte(val))
}

func (h *Handler) AuthPublic(c *gin.Context) {
	var body struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	token, code, err := h.svc.AuthenticatePublic(c.Param("slug"), body.Password)
	if err != nil {
		logger.WithCtx(c).Error("share.auth_public failed", zap.Error(err), zap.Int("code", code), zap.String("slug", c.Param("slug")))
		switch code {
		case 401:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_password"})
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		case 410:
			c.JSON(http.StatusGone, gin.H{"error": "expired"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) GetSettings(c *gin.Context) {
	uid, ok := middleware.UID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	sl, code, err := h.svc.GetSettingsForUser(uid, uint(id))
	if err != nil {
		logger.WithCtx(c).Error("share.get_settings failed", zap.Error(err), zap.Int("code", code), zap.String("id", c.Param("id")))
		switch code {
		case 403:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		}
		return
	}
	views := sl.Views
	lastAccessAt := sl.LastAccessAt
	if cache.RDB != nil {
		if v := cache.RDB.Get(c, common.RedisKeyViews.F(sl.Slug)).Val(); v != "" {
			if n, err := strconv.ParseUint(v, 10, 64); err == nil {
				views = n
			}
		}
		if v := cache.RDB.Get(c, common.RedisKeyShareLastAccess.F(sl.Slug)).Val(); v != "" {
			if ms, err := strconv.ParseInt(v, 10, 64); err == nil && ms > 0 {
				t := time.UnixMilli(ms)
				lastAccessAt = &t
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"slug":         sl.Slug,
		"isPublic":     sl.IsPublic,
		"hasPassword":  sl.Password != "",
		"password":     sl.Password,
		"expiresAt":    sl.ExpiresAt,
		"views":        views,
		"lastAccessAt": lastAccessAt,
		"url":          "/public/resumes/" + sl.Slug,
		"apiUrl":       "/public/resumes/" + sl.Slug,
		"pageUrl":      "/#/public/" + sl.Slug,
	})
}

func (h *Handler) UpdateSettings(c *gin.Context) {
	uid, ok := middleware.UID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	var body struct {
		IsPublic  *bool   `json:"isPublic"`
		Password  *string `json:"password"`
		ExpiresAt *string `json:"expiresAt"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.WithCtx(c).Error("share.update_settings bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	in := UpdateSettingsInput{IsPublic: body.IsPublic, Password: body.Password}
	if body.ExpiresAt != nil {
		v := strings.TrimSpace(*body.ExpiresAt)
		if v == "" {
			in.ClearExpiresAt = true
		} else {
			tm, err := time.Parse(time.RFC3339, v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expiresAt"})
				return
			}
			in.ExpiresAt = &tm
		}
	}
	sl, code, err := h.svc.UpdateSettingsForUser(uid, uint(id), in)
	if err != nil {
		logger.WithCtx(c).Error("share.update_settings failed", zap.Error(err), zap.Int("code", code), zap.String("id", c.Param("id")))
		switch code {
		case 403:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"slug":        sl.Slug,
		"isPublic":    sl.IsPublic,
		"hasPassword": sl.Password != "",
		"password":    sl.Password,
		"expiresAt":   sl.ExpiresAt,
	})
}
