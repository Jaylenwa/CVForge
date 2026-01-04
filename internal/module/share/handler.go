package share

import (
	"net/http"

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
	sl, code, err := h.svc.PublishResumeForUser(uid, c.Param("id"))
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
	val, code, err := h.svc.GetPublicPayload(c.Param("slug"))
	if err != nil {
		logger.WithCtx(c).Error("share.get_public failed", zap.Error(err), zap.Int("code", code), zap.String("slug", c.Param("slug")))
		switch code {
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		}
		return
	}
	c.Data(http.StatusOK, "application/json", []byte(val))
}
