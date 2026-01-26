package library

import (
	"net/http"
	"strconv"
	"strings"

	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

type GenerateMode string

const (
	GenerateModeSkip   GenerateMode = "skip"
	GenerateModeUpdate GenerateMode = "update"
)

type GenerateVariantsReq struct {
	RoleID            uint     `json:"roleId"`
	PresetID          uint     `json:"presetId"`
	LayoutTemplateIDs []string `json:"layoutTemplateIds"`
	NamePrefix        string   `json:"namePrefix"`
	Tags              string   `json:"tags"`
	IsPremium         *bool    `json:"isPremium"`
	IsActive          *bool    `json:"isActive"`
	Mode              string   `json:"mode"`
}

type GenerateVariantsResult struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
	Skipped int `json:"skipped"`
	Failed  int `json:"failed"`
	Items   []struct {
		LayoutTemplateID string `json:"layoutTemplateId"`
		VariantID        uint   `json:"variantId"`
		Action           string `json:"action"`
		Error            string `json:"error,omitempty"`
	} `json:"items"`
}

func (h *AdminHandler) AdminListVariants(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("pageSize"), 20)
	q := c.Query("q")
	role := c.Query("roleId")
	category := c.Query("categoryId")
	template := c.Query("templateId")
	items, total, err := h.svc.repo.AdminListTemplateVariants(page, size, q, role, category, template)
	if err != nil {
		logger.WithCtx(c).Error("library.admin.variants.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, AdminPageResp[TemplateVariant]{Items: items, Page: clampPage(page), PageSize: clampPageSize(size), Total: total})
}

func (h *AdminHandler) AdminCreateVariant(c *gin.Context) {
	var body struct {
		Name                     string `json:"name"`
		LayoutTemplateExternalID string `json:"layoutTemplateExternalId"`
		PresetID                 uint   `json:"presetId"`
		RoleID                   uint   `json:"roleId"`
		Tags                     string `json:"tags"`
		UsageCount               *int   `json:"usageCount"`
		IsPremium                *bool  `json:"isPremium"`
		IsActive                 *bool  `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.Name) == "" || body.PresetID == 0 || body.RoleID == 0 || strings.TrimSpace(body.LayoutTemplateExternalID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	m := TemplateVariant{
		Name:                     strings.TrimSpace(body.Name),
		LayoutTemplateExternalID: strings.TrimSpace(body.LayoutTemplateExternalID),
		PresetID:                 body.PresetID,
		RoleID:                   body.RoleID,
		Tags:                     strings.TrimSpace(body.Tags),
	}
	if body.UsageCount != nil {
		m.UsageCount = *body.UsageCount
	}
	if body.IsPremium != nil {
		m.IsPremium = *body.IsPremium
	}
	if body.IsActive != nil {
		m.IsActive = *body.IsActive
	}
	if err := h.svc.repo.AdminCreateTemplateVariant(&m); err != nil {
		logger.WithCtx(c).Error("library.admin.variants.create failed", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminPatchVariant(c *gin.Context) {
	var body struct {
		Name                     *string `json:"name"`
		LayoutTemplateExternalID *string `json:"layoutTemplateExternalId"`
		PresetID                 *uint   `json:"presetId"`
		RoleID                   *uint   `json:"roleId"`
		Tags                     *string `json:"tags"`
		UsageCount               *int    `json:"usageCount"`
		IsPremium                *bool   `json:"isPremium"`
		IsActive                 *bool   `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	patch := map[string]any{}
	if body.Name != nil {
		patch["name"] = strings.TrimSpace(*body.Name)
	}
	if body.LayoutTemplateExternalID != nil {
		patch["layout_template_external_id"] = strings.TrimSpace(*body.LayoutTemplateExternalID)
	}
	if body.PresetID != nil {
		patch["preset_id"] = *body.PresetID
	}
	if body.RoleID != nil {
		patch["role_id"] = *body.RoleID
	}
	if body.Tags != nil {
		patch["tags"] = strings.TrimSpace(*body.Tags)
	}
	if body.UsageCount != nil {
		patch["usage_count"] = *body.UsageCount
	}
	if body.IsPremium != nil {
		patch["is_premium"] = *body.IsPremium
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
	if err := h.svc.repo.AdminPatchTemplateVariant(uint(id), patch); err != nil {
		logger.WithCtx(c).Error("library.admin.variants.patch failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminDeleteVariant(c *gin.Context) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := h.svc.repo.AdminDeleteTemplateVariant(uint(id)); err != nil {
		logger.WithCtx(c).Error("library.admin.variants.delete failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminGenerateVariants(c *gin.Context) {
	var req GenerateVariantsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	mode := GenerateMode(strings.TrimSpace(req.Mode))
	if mode != GenerateModeSkip && mode != GenerateModeUpdate {
		mode = GenerateModeSkip
	}
	roleID := req.RoleID
	presetID := req.PresetID
	if roleID == 0 || presetID == 0 || len(req.LayoutTemplateIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	isPremium := false
	if req.IsPremium != nil {
		isPremium = *req.IsPremium
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	tags := strings.TrimSpace(req.Tags)
	namePrefix := strings.TrimSpace(req.NamePrefix)

	var out GenerateVariantsResult
	err := h.svc.repo.db.Transaction(func(tx *gorm.DB) error {
		for _, tplID := range req.LayoutTemplateIDs {
			tplID = strings.TrimSpace(tplID)
			if tplID == "" {
				continue
			}
			name := namePrefix
			if name == "" {
				tplName, _ := h.svc.repo.getTemplateName(tplID)
				roleName, _ := h.svc.repo.getRoleName(roleID)
				presetName, _ := h.svc.repo.getPresetName(presetID)
				name = strings.TrimSpace(roleName + " " + tplName + " " + presetName)
			}
			v := TemplateVariant{
				Name:                     name,
				LayoutTemplateExternalID: tplID,
				PresetID:                 presetID,
				RoleID:                   roleID,
				Tags:                     tags,
				IsPremium:                isPremium,
				IsActive:                 isActive,
			}
			action, genErr := h.svc.repo.upsertVariant(tx, &v, mode)
			item := struct {
				LayoutTemplateID string `json:"layoutTemplateId"`
				VariantID        uint   `json:"variantId"`
				Action           string `json:"action"`
				Error            string `json:"error,omitempty"`
			}{
				LayoutTemplateID: tplID,
				Action:           action,
			}
			if got, e := h.svc.repo.findVariantByCombo(roleID, presetID, tplID); e == nil {
				item.VariantID = got.ID
			}
			if genErr != nil {
				out.Failed++
				item.Action = "failed"
				item.Error = genErr.Error()
			} else {
				switch action {
				case "created":
					out.Created++
				case "updated":
					out.Updated++
				case "skipped":
					out.Skipped++
				}
			}
			out.Items = append(out.Items, item)
		}
		return nil
	})
	if err != nil {
		logger.WithCtx(c).Error("library.admin.variants.generate failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, out)
}
