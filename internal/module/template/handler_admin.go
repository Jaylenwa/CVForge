package template

import (
	"net/http"
	"strings"

	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdminHandler struct {
	svc *Service
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{svc: NewService()}
}

func (h *AdminHandler) AdminCreate(c *gin.Context) {
	var body struct {
		ExternalID string `json:"externalId"`
		Name       string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.ExternalID == "" || body.Name == "" {
		logger.WithCtx(c).Error("template.admin_create bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	mt := Template{
		ExternalID: body.ExternalID,
		Name:       body.Name,
	}
	if err := h.svc.Create(mt); err != nil {
		logger.WithCtx(c).Error("template.admin_create conflict", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminPatch(c *gin.Context) {
	var body struct {
		Name *string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.WithCtx(c).Error("template.admin_patch bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	err := h.svc.Update(c.Param("id"), func(t *Template) {
		if body.Name != nil {
			t.Name = strings.TrimSpace(*body.Name)
		}
	})
	if err != nil {
		logger.WithCtx(c).Error("template.admin_patch failed", zap.Error(err), zap.String("id", c.Param("id")))
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminDelete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		logger.WithCtx(c).Error("template.admin_delete failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminSeed(c *gin.Context) {
	var body []struct {
		ExternalID string `json:"externalId"`
		Name       string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.WithCtx(c).Error("template.admin_seed bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	var templates []Template
	for _, item := range body {
		if item.ExternalID == "" || item.Name == "" {
			continue
		}
		templates = append(templates, Template{
			ExternalID: item.ExternalID,
			Name:       item.Name,
		})
	}

	if err := h.svc.Seed(templates); err != nil {
		logger.WithCtx(c).Error("template.admin_seed failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
