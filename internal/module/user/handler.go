package user

import (
	"net/http"

	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	svc *Service
}

func NewHandler() *Handler {
	return &Handler{svc: NewService()}
}

func (h *Handler) Me(c *gin.Context) {
	uidVal, _ := c.Get("uid")
	u, err := h.svc.GetMe(uidVal)
	if err != nil {
		logger.WithCtx(c).Error("user.me failed", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	email := ""
	if u.Email != nil {
		email = *u.Email
	}
	c.JSON(http.StatusOK, gin.H{
		"email":     email,
		"name":      u.Name,
		"avatarUrl": u.AvatarURL,
		"language":  u.Language,
		"role":      u.Role,
	})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	uidVal, _ := c.Get("uid")
	var body struct{ Name, AvatarURL, Language string }
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.WithCtx(c).Error("user.update_profile bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := h.svc.UpdateProfile(uidVal, body.Name, body.AvatarURL, body.Language); err != nil {
		logger.WithCtx(c).Error("user.update_profile failed", zap.Error(err))
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) UpdatePassword(c *gin.Context) {
	uidVal, _ := c.Get("uid")
	var body struct{ CurrentPassword, NewPassword string }
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.WithCtx(c).Error("user.update_password bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := h.svc.UpdatePassword(uidVal, body.CurrentPassword, body.NewPassword); err != nil {
		logger.WithCtx(c).Error("user.update_password failed", zap.Error(err))
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if err == gorm.ErrInvalidData {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid current password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
