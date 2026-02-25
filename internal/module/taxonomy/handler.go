package taxonomy

import (
	"net/http"
	"strconv"

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
	language := normalizeLanguage(c.Query("language"))
	items, err := h.svc.ListJobCategories(language)
	if err != nil {
		logger.WithCtx(c).Error("taxonomy.categories.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]JobCategoryDTO, 0, len(items))
	for _, it := range items {
		out = append(out, JobCategoryDTO{
			ID:       it.ID,
			Name:     it.Name,
			ParentID: it.ParentID,
			OrderNum: it.OrderNum,
			IsActive: it.IsActive,
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": out})
}

func (h *Handler) ListRoles(c *gin.Context) {
	language := normalizeLanguage(c.Query("language"))
	var categoryID uint
	if v := c.Query("categoryId"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			categoryID = uint(n)
		}
	}
	q := c.Query("q")
	items, err := h.svc.ListJobRoles(language, categoryID, q)
	if err != nil {
		logger.WithCtx(c).Error("taxonomy.roles.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]JobRoleDTO, 0, len(items))
	for _, it := range items {
		out = append(out, JobRoleDTO{
			ID:         it.ID,
			CategoryID: it.CategoryID,
			Name:       it.Name,
			OrderNum:   it.OrderNum,
			IsActive:   it.IsActive,
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": out})
}
