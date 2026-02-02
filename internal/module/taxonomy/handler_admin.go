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

func (h *AdminHandler) AdminListCategories(c *gin.Context) {
	language := normalizeLanguage(c.Query("language"))
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("pageSize"), 20)
	q := c.Query("q")
	var parentID *uint
	if v := strings.TrimSpace(c.Query("parentId")); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			x := uint(n)
			parentID = &x
		}
	}
	items, total, err := h.svc.repo.AdminListJobCategories(page, size, q, parentID)
	if err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.categories.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	ids := make([]uint, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.ID)
	}
	i18nList, err := h.svc.repo.ListJobCategoryI18n(ids)
	if err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.categories.i18n.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	namesByID := map[uint]map[string]string{}
	for _, it := range i18nList {
		m := namesByID[it.JobCategoryID]
		if m == nil {
			m = map[string]string{}
			namesByID[it.JobCategoryID] = m
		}
		m[it.Language] = it.Name
	}
	out := make([]AdminJobCategoryDTO, 0, len(items))
	for _, it := range items {
		names := namesByID[it.ID]
		out = append(out, AdminJobCategoryDTO{
			ID:       it.ID,
			Name:     pickName(names, language),
			Names:    names,
			ParentID: it.ParentID,
			OrderNum: it.OrderNum,
			IsActive: it.IsActive,
		})
	}
	c.JSON(http.StatusOK, AdminPageResp[AdminJobCategoryDTO]{Items: out, Page: clampPage(page), PageSize: clampPageSize(size), Total: total})
}

func (h *AdminHandler) AdminCreateCategory(c *gin.Context) {
	var body struct {
		Name     string            `json:"name"`
		Names    map[string]string `json:"names"`
		ParentID *uint             `json:"parentId"`
		OrderNum *int              `json:"orderNum"`
		IsActive *bool             `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	m := JobCategory{
		ParentID: body.ParentID,
	}
	if body.OrderNum != nil {
		m.OrderNum = *body.OrderNum
	}
	if body.IsActive != nil {
		m.IsActive = *body.IsActive
	}
	if err := h.svc.repo.AdminCreateJobCategoryWithNames(&m, names); err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.categories.create failed", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminPatchCategory(c *gin.Context) {
	var body struct {
		Name     *string           `json:"name"`
		Names    map[string]string `json:"names"`
		ParentID *uint             `json:"parentId"`
		OrderNum *int              `json:"orderNum"`
		IsActive *bool             `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	patch := map[string]any{}
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
	if body.ParentID != nil {
		patch["parent_id"] = *body.ParentID
	}
	if body.OrderNum != nil {
		patch["order_num"] = *body.OrderNum
	}
	if body.IsActive != nil {
		patch["is_active"] = *body.IsActive
	}
	if len(patch) == 0 && len(names) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := h.svc.repo.AdminPatchJobCategoryWithNames(uint(id), patch, names); err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.categories.patch failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminDeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := h.svc.repo.AdminDeleteJobCategory(uint(id)); err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.categories.delete failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminListRoles(c *gin.Context) {
	language := normalizeLanguage(c.Query("language"))
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("pageSize"), 20)
	q := c.Query("q")
	var categoryID *uint
	if v := strings.TrimSpace(c.Query("categoryId")); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			x := uint(n)
			categoryID = &x
		}
	}
	items, total, err := h.svc.repo.AdminListJobRoles(page, size, q, categoryID)
	if err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.roles.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	ids := make([]uint, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.ID)
	}
	i18nList, err := h.svc.repo.ListJobRoleI18n(ids)
	if err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.roles.i18n.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	namesByID := map[uint]map[string]string{}
	for _, it := range i18nList {
		m := namesByID[it.JobRoleID]
		if m == nil {
			m = map[string]string{}
			namesByID[it.JobRoleID] = m
		}
		m[it.Language] = it.Name
	}
	out := make([]AdminJobRoleDTO, 0, len(items))
	for _, it := range items {
		names := namesByID[it.ID]
		out = append(out, AdminJobRoleDTO{
			ID:         it.ID,
			CategoryID: it.CategoryID,
			Name:       pickName(names, language),
			Names:      names,
			OrderNum:   it.OrderNum,
			IsActive:   it.IsActive,
		})
	}
	c.JSON(http.StatusOK, AdminPageResp[AdminJobRoleDTO]{Items: out, Page: clampPage(page), PageSize: clampPageSize(size), Total: total})
}

func (h *AdminHandler) AdminCreateRole(c *gin.Context) {
	var body struct {
		CategoryID uint              `json:"categoryId"`
		Name       string            `json:"name"`
		Names      map[string]string `json:"names"`
		OrderNum   *int              `json:"orderNum"`
		IsActive   *bool             `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.CategoryID == 0 {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	m := JobRole{
		CategoryID: body.CategoryID,
	}
	if body.OrderNum != nil {
		m.OrderNum = *body.OrderNum
	}
	if body.IsActive != nil {
		m.IsActive = *body.IsActive
	}
	if err := h.svc.repo.AdminCreateJobRoleWithNames(&m, names); err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.roles.create failed", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminPatchRole(c *gin.Context) {
	var body struct {
		CategoryID *uint             `json:"categoryId"`
		Name       *string           `json:"name"`
		Names      map[string]string `json:"names"`
		OrderNum   *int              `json:"orderNum"`
		IsActive   *bool             `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	patch := map[string]any{}
	names := body.Names
	if names == nil {
		names = map[string]string{}
	}
	if body.CategoryID != nil {
		patch["category_id"] = *body.CategoryID
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
	if body.OrderNum != nil {
		patch["order_num"] = *body.OrderNum
	}
	if body.IsActive != nil {
		patch["is_active"] = *body.IsActive
	}
	if len(patch) == 0 && len(names) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := h.svc.repo.AdminPatchJobRoleWithNames(uint(id), patch, names); err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.roles.patch failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminDeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := h.svc.repo.AdminDeleteJobRole(uint(id)); err != nil {
		logger.WithCtx(c).Error("taxonomy.admin.roles.delete failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
