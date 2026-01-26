package library

import (
	"net/http"
	"strconv"
	"strings"

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
	var roleID uint
	if v := strings.TrimSpace(c.Query("roleId")); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			roleID = uint(n)
		}
	}
	var categoryID uint
	if v := strings.TrimSpace(c.Query("categoryId")); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			categoryID = uint(n)
		}
	}
	q := c.Query("q")
	items, err := h.svc.ListTemplateVariants(roleID, categoryID, q)
	if err != nil {
		logger.WithCtx(c).Error("library.variants.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	out := make([]TemplateVariantDTO, 0, len(items))
	for _, v := range items {
		out = append(out, TemplateVariantDTO{
			ID:                       v.ID,
			Name:                     v.Name,
			LayoutTemplateExternalID: v.LayoutTemplateExternalID,
			PresetID:                 v.PresetID,
			RoleID:                   v.RoleID,
			Tags:                     splitTags(v.Tags),
			UsageCount:               v.UsageCount,
			IsPremium:                v.IsPremium,
			IsActive:                 v.IsActive,
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": out})
}
