package library

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

func (h *Handler) ListVariants(c *gin.Context) {
	role := c.Query("roleId")
	category := c.Query("categoryId")
	q := c.Query("q")
	items, err := h.svc.ListTemplateVariants(role, category, q)
	if err != nil {
		logger.WithCtx(c).Error("library.variants.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	out := make([]TemplateVariantDTO, 0, len(items))
	for _, v := range items {
		out = append(out, TemplateVariantDTO{
			ExternalID:               v.ExternalID,
			Name:                     v.Name,
			LayoutTemplateExternalID: v.LayoutTemplateExternalID,
			PresetExternalID:         v.PresetExternalID,
			RoleExternalID:           v.RoleExternalID,
			Tags:                     splitTags(v.Tags),
			UsageCount:               v.UsageCount,
			IsPremium:                v.IsPremium,
			IsActive:                 v.IsActive,
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": out})
}

