package preset

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

func (h *AdminHandler) AdminListPresets(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("pageSize"), 20)
	q := c.Query("q")
	role := c.Query("roleExternalId")
	lang := c.Query("language")
	items, total, err := h.svc.repo.AdminListContentPresets(page, size, q, role, lang)
	if err != nil {
		logger.WithCtx(c).Error("preset.admin.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, AdminPageResp[ContentPreset]{Items: items, Page: clampPage(page), PageSize: clampPageSize(size), Total: total})
}

func (h *AdminHandler) AdminCreatePreset(c *gin.Context) {
	var body struct {
		ExternalID     string `json:"externalId"`
		Name           string `json:"name"`
		Language       string `json:"language"`
		RoleExternalID string `json:"roleExternalId"`
		Tags           string `json:"tags"`
		DataJSON       string `json:"dataJson"`
		IsActive       *bool  `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.ExternalID) == "" || strings.TrimSpace(body.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	m := ContentPreset{
		ExternalID:     strings.TrimSpace(body.ExternalID),
		Name:           strings.TrimSpace(body.Name),
		Language:       strings.TrimSpace(body.Language),
		RoleExternalID: strings.TrimSpace(body.RoleExternalID),
		Tags:           strings.TrimSpace(body.Tags),
		DataJSON:       strings.TrimSpace(body.DataJSON),
	}
	if body.IsActive != nil {
		m.IsActive = *body.IsActive
	}
	if err := h.svc.repo.AdminCreateContentPreset(&m); err != nil {
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

func (h *AdminHandler) AdminPatchPreset(c *gin.Context) {
	var body struct {
		Name           *string `json:"name"`
		Language       *string `json:"language"`
		RoleExternalID *string `json:"roleExternalId"`
		Tags           *string `json:"tags"`
		DataJSON       *string `json:"dataJson"`
		IsActive       *bool   `json:"isActive"`
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
	if body.RoleExternalID != nil {
		patch["role_external_id"] = strings.TrimSpace(*body.RoleExternalID)
	}
	if body.Tags != nil {
		patch["tags"] = strings.TrimSpace(*body.Tags)
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
	if err := h.svc.repo.AdminPatchContentPreset(c.Param("id"), patch); err != nil {
		logger.WithCtx(c).Error("preset.admin.patch failed", zap.Error(err))
		if strings.HasPrefix(err.Error(), "invalid_") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid preset"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminDeletePreset(c *gin.Context) {
	if err := h.svc.repo.AdminDeleteContentPreset(c.Param("id")); err != nil {
		logger.WithCtx(c).Error("preset.admin.delete failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
