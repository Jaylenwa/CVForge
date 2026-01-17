package catalog

import (
	"net/http"

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

func (h *Handler) ListJobCategories(c *gin.Context) {
	items, err := h.svc.ListJobCategories()
	if err != nil {
		logger.WithCtx(c).Error("catalog.job_categories.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) ListJobRoles(c *gin.Context) {
	category := c.Query("category")
	q := c.Query("q")
	items, err := h.svc.ListJobRoles(category, q)
	if err != nil {
		logger.WithCtx(c).Error("catalog.job_roles.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) ListTemplateVariants(c *gin.Context) {
	role := c.Query("role")
	category := c.Query("category")
	q := c.Query("q")
	items, err := h.svc.ListTemplateVariants(role, category, q)
	if err != nil {
		logger.WithCtx(c).Error("catalog.template_variants.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) GetContentPreset(c *gin.Context) {
	p, err := h.svc.GetContentPresetByExternal(c.Param("id"))
	if err != nil {
		logger.WithCtx(c).Error("catalog.content_preset.get failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

