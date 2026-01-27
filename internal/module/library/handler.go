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

func (h *Handler) ListTemplates(c *gin.Context) {
	var roleID uint
	if v := strings.TrimSpace(c.Query("roleId")); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			roleID = uint(n)
		}
	}
	items, err := h.svc.ListTemplateLibraryItems(roleID)
	if err != nil {
		logger.WithCtx(c).Error("library.templates.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	out := make([]TemplateLibraryItemDTO, 0, len(items))
	for _, it := range items {
		out = append(out, TemplateLibraryItemDTO{
			TemplateExternalID: it.TemplateExternalID,
			Name:               it.Name,
			Tags:               splitTags(it.Tags),
			UsageCount:         it.UsageCount,
			GlobalUsageCount:   it.GlobalUsageCount,
			PresetID:           it.PresetID,
			RoleID:             it.RoleID,
			IsPremium:          it.IsPremium,
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": out})
}
