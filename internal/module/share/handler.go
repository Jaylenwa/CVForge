package share

import (
	"net/http"

	"openresume/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	svc *Service
}

func NewHandler(db *gorm.DB, rdb *redis.Client) *Handler {
	return &Handler{svc: NewService(db, rdb)}
}

func (h *Handler) PublishResume(c *gin.Context) {
	uid, ok := middleware.UID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	sl, code, err := h.svc.PublishResumeForUser(uid, c.Param("id"))
	if err != nil {
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
	c.JSON(http.StatusOK, gin.H{"slug": sl.Slug, "url": "/public/resumes/" + sl.Slug})
}

func (h *Handler) GetPublic(c *gin.Context) {
	val, code, err := h.svc.GetPublicPayload(c.Param("slug"))
	if err != nil {
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
