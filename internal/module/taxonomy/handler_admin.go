package taxonomy

import (
	"net/http"
	"strconv"
	"strings"

	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminHandler struct {
	svc *Service
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{svc: NewService()}
}

type AdminPageResp[T any] struct {
	Items    []T   `json:"items"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

func parseIntDefault(s string, d int) int {
	if s = strings.TrimSpace(s); s == "" {
		return d
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return n
}

func clampPage(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

func clampPageSize(size int) int {
	if size <= 0 {
		return 20
	}
	if size > 100 {
		return 100
	}
	return size
}

func (h *AdminHandler) AdminListCategories(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("pageSize"), 20)
	q := c.Query("q")
	parent := c.Query("parent")
	items, total, err := h.svc.repo.AdminListJobCategories(page, size, q, parent)
	if err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.categories.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, AdminPageResp[JobCategory]{Items: items, Page: clampPage(page), PageSize: clampPageSize(size), Total: total})
}

func (h *AdminHandler) AdminCreateCategory(c *gin.Context) {
	var body struct {
		ExternalID       string `json:"externalId"`
		Name             string `json:"name"`
		ParentExternalID string `json:"parentExternalId"`
		OrderNum         *int   `json:"orderNum"`
		IsActive         *bool  `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.ExternalID) == "" || strings.TrimSpace(body.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	m := JobCategory{
		ExternalID:       strings.TrimSpace(body.ExternalID),
		Name:             strings.TrimSpace(body.Name),
		ParentExternalID: strings.TrimSpace(body.ParentExternalID),
	}
	if body.OrderNum != nil {
		m.OrderNum = *body.OrderNum
	}
	if body.IsActive != nil {
		m.IsActive = *body.IsActive
	}
	if err := h.svc.repo.AdminCreateJobCategory(&m); err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.categories.create failed", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminPatchCategory(c *gin.Context) {
	var body struct {
		Name             *string `json:"name"`
		ParentExternalID *string `json:"parentExternalId"`
		OrderNum         *int    `json:"orderNum"`
		IsActive         *bool   `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	patch := map[string]any{}
	if body.Name != nil {
		patch["name"] = strings.TrimSpace(*body.Name)
	}
	if body.ParentExternalID != nil {
		patch["parent_external_id"] = strings.TrimSpace(*body.ParentExternalID)
	}
	if body.OrderNum != nil {
		patch["order_num"] = *body.OrderNum
	}
	if body.IsActive != nil {
		patch["is_active"] = *body.IsActive
	}
	if len(patch) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}
	if err := h.svc.repo.AdminPatchJobCategory(c.Param("id"), patch); err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.categories.patch failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminDeleteCategory(c *gin.Context) {
	if err := h.svc.repo.AdminDeleteJobCategory(c.Param("id")); err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.categories.delete failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminListRoles(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("pageSize"), 20)
	q := c.Query("q")
	category := c.Query("categoryId")
	items, total, err := h.svc.repo.AdminListJobRoles(page, size, q, category)
	if err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.roles.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, AdminPageResp[JobRole]{Items: items, Page: clampPage(page), PageSize: clampPageSize(size), Total: total})
}

func (h *AdminHandler) AdminCreateRole(c *gin.Context) {
	var body struct {
		ExternalID         string `json:"externalId"`
		CategoryExternalID string `json:"categoryExternalId"`
		Name               string `json:"name"`
		Tags               string `json:"tags"`
		OrderNum           *int   `json:"orderNum"`
		IsActive           *bool  `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.ExternalID) == "" || strings.TrimSpace(body.Name) == "" || strings.TrimSpace(body.CategoryExternalID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	m := JobRole{
		ExternalID:         strings.TrimSpace(body.ExternalID),
		CategoryExternalID: strings.TrimSpace(body.CategoryExternalID),
		Name:               strings.TrimSpace(body.Name),
		Tags:               strings.TrimSpace(body.Tags),
	}
	if body.OrderNum != nil {
		m.OrderNum = *body.OrderNum
	}
	if body.IsActive != nil {
		m.IsActive = *body.IsActive
	}
	if err := h.svc.repo.AdminCreateJobRole(&m); err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.roles.create failed", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminPatchRole(c *gin.Context) {
	var body struct {
		CategoryExternalID *string `json:"categoryExternalId"`
		Name               *string `json:"name"`
		Tags               *string `json:"tags"`
		OrderNum           *int    `json:"orderNum"`
		IsActive           *bool   `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	patch := map[string]any{}
	if body.CategoryExternalID != nil {
		patch["category_external_id"] = strings.TrimSpace(*body.CategoryExternalID)
	}
	if body.Name != nil {
		patch["name"] = strings.TrimSpace(*body.Name)
	}
	if body.Tags != nil {
		patch["tags"] = strings.TrimSpace(*body.Tags)
	}
	if body.OrderNum != nil {
		patch["order_num"] = *body.OrderNum
	}
	if body.IsActive != nil {
		patch["is_active"] = *body.IsActive
	}
	if len(patch) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}
	if err := h.svc.repo.AdminPatchJobRole(c.Param("id"), patch); err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.roles.patch failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminDeleteRole(c *gin.Context) {
	if err := h.svc.repo.AdminDeleteJobRole(c.Param("id")); err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.roles.delete failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

