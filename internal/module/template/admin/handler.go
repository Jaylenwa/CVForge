package admin

import (
	"net/http"
	"strings"

	"cvforge/internal/pkg/logger"

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

func pickName(names map[string]string, language string) string {
	if len(names) == 0 {
		return ""
	}
	if v := strings.TrimSpace(names[language]); v != "" {
		return v
	}
	if v := strings.TrimSpace(names["zh"]); v != "" {
		return v
	}
	for _, v := range names {
		if v = strings.TrimSpace(v); v != "" {
			return v
		}
	}
	return ""
}

func (h *Handler) AdminCreate(c *gin.Context) {
	var body struct {
		ExternalID string            `json:"externalId"`
		Name       string            `json:"name"`
		Names      map[string]string `json:"names"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.ExternalID) == "" {
		logger.WithCtx(c).Error("template.admin_create bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	names := body.Names
	if names == nil {
		names = map[string]string{}
	}
	if strings.TrimSpace(names["zh"]) == "" && strings.TrimSpace(body.Name) != "" {
		names["zh"] = body.Name
	}
	if strings.TrimSpace(pickName(names, "zh")) == "" && strings.TrimSpace(pickName(names, "en")) == "" {
		logger.WithCtx(c).Error("template.admin_create bad request empty names")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if err := h.svc.Create(body.ExternalID, names); err != nil {
		logger.WithCtx(c).Error("template.admin_create conflict", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) AdminPatch(c *gin.Context) {
	var body struct {
		Name  *string           `json:"name"`
		Names map[string]string `json:"names"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.WithCtx(c).Error("template.admin_patch bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	names := body.Names
	if names == nil {
		names = map[string]string{}
	}
	if body.Name != nil {
		trimmed := strings.TrimSpace(*body.Name)
		if trimmed != "" {
			names["zh"] = trimmed
		}
	}
	if v, ok := names["zh"]; ok {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			names["zh"] = trimmed
		} else {
			delete(names, "zh")
		}
	}
	if v, ok := names["en"]; ok {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			names["en"] = trimmed
		} else {
			delete(names, "en")
		}
	}

	if len(names) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}

	err := h.svc.UpdateNames(c.Param("id"), names)
	if err != nil {
		logger.WithCtx(c).Error("template.admin_patch failed", zap.Error(err), zap.String("id", c.Param("id")))
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) AdminDelete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		logger.WithCtx(c).Error("template.admin_delete failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

