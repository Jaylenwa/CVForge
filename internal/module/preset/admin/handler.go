package admin

import (
	"net/http"
	"strconv"
	"strings"

	presetmod "cvforge/internal/module/preset"
	"cvforge/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	svc *Service
}

func NewHandler() *Handler {
	return &Handler{svc: NewService()}
}

type PageResp[T any] struct {
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

func (h *Handler) AdminListPresets(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("pageSize"), 20)
	q := c.Query("q")
	role := c.Query("roleId")
	lang := c.Query("language")
	items, total, err := h.svc.AdminListContentPresets(page, size, q, role, lang)
	if err != nil {
		logger.WithCtx(c).Error("preset.admin.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, PageResp[presetmod.ContentPreset]{Items: items, Page: clampPage(page), PageSize: clampPageSize(size), Total: total})
}

func (h *Handler) AdminCreatePreset(c *gin.Context) {
	var body struct {
		Name     string `json:"name"`
		Language string `json:"language"`
		RoleID   uint   `json:"roleId"`
		DataJSON string `json:"dataJson"`
		IsActive *bool  `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.Name) == "" || body.RoleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	m := presetmod.ContentPreset{
		Name:     strings.TrimSpace(body.Name),
		Language: strings.TrimSpace(body.Language),
		RoleID:   body.RoleID,
		DataJSON: strings.TrimSpace(body.DataJSON),
		IsActive: true,
	}
	if body.IsActive != nil {
		m.IsActive = *body.IsActive
	}
	if err := h.svc.AdminCreateContentPreset(&m); err != nil {
		logger.WithCtx(c).Error("preset.admin.create failed", zap.Error(err))
		if strings.HasPrefix(err.Error(), "invalid_") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid preset"})
			return
		}
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) AdminPatchPreset(c *gin.Context) {
	var body struct {
		Name     *string `json:"name"`
		Language *string `json:"language"`
		RoleID   *uint   `json:"roleId"`
		DataJSON *string `json:"dataJson"`
		IsActive *bool   `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	patch := map[string]any{}
	if body.Name != nil {
		patch["name"] = strings.TrimSpace(*body.Name)
	}
	if body.Language != nil {
		patch["language"] = strings.TrimSpace(*body.Language)
	}
	if body.RoleID != nil {
		patch["role_id"] = *body.RoleID
	}
	if body.DataJSON != nil {
		patch["data_json"] = strings.TrimSpace(*body.DataJSON)
	}
	if body.IsActive != nil {
		patch["is_active"] = *body.IsActive
	}
	if len(patch) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := h.svc.AdminPatchContentPreset(uint(id), patch); err != nil {
		logger.WithCtx(c).Error("preset.admin.patch failed", zap.Error(err))
		if strings.HasPrefix(err.Error(), "invalid_") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid preset"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) AdminDeletePreset(c *gin.Context) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := h.svc.AdminDeleteContentPreset(uint(id)); err != nil {
		logger.WithCtx(c).Error("preset.admin.delete failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

