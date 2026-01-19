package taxonomy

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

func (h *Handler) ListCategories(c *gin.Context) {
	items, err := h.svc.ListJobCategories()
	if err != nil {
		logger.WithCtx(c).Error("taxonomy.categories.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	out := make([]JobCategoryDTO, 0, len(items))
	for _, it := range items {
		out = append(out, JobCategoryDTO{
			ExternalID:       it.ExternalID,
			Name:             it.Name,
			ParentExternalID: it.ParentExternalID,
			OrderNum:         it.OrderNum,
			IsActive:         it.IsActive,
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": out})
}

func (h *Handler) ListRoles(c *gin.Context) {
	category := c.Query("categoryId")
	q := c.Query("q")
	items, err := h.svc.ListJobRoles(category, q)
	if err != nil {
		logger.WithCtx(c).Error("taxonomy.roles.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	out := make([]JobRoleDTO, 0, len(items))
	for _, it := range items {
		out = append(out, JobRoleDTO{
			ExternalID:         it.ExternalID,
			CategoryExternalID: it.CategoryExternalID,
			Name:               it.Name,
			Tags:               splitTags(it.Tags),
			OrderNum:           it.OrderNum,
			IsActive:           it.IsActive,
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": out})
}

