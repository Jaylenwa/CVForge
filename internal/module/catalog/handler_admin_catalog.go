package catalog

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"

	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func parseIntDefault(s string, d int) int {
	if s == "" {
		return d
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return d
	}
	return n
}

func sanitizeExternalID(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	out := b.String()
	out = strings.Trim(out, "_-")
	for strings.Contains(out, "__") {
		out = strings.ReplaceAll(out, "__", "_")
	}
	return out
}

func stableExternalID(parts ...string) string {
	raw := strings.Join(parts, "_")
	out := sanitizeExternalID(raw)
	if len(out) <= 64 {
		return out
	}
	h := sha1.Sum([]byte(out))
	suffix := hex.EncodeToString(h[:])[:8]
	if len(out) > 55 {
		out = out[:55]
	}
	return strings.TrimRight(out, "_-") + "_" + suffix
}

func (h *AdminHandler) AdminListJobCategories(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("pageSize"), 20)
	q := c.Query("q")
	parent := c.Query("parent")
	items, total, err := h.svc.repo.AdminListJobCategories(page, size, q, parent)
	if err != nil {
		logger.WithCtx(c).Error("catalog.admin.job_categories.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, AdminPageResp[JobCategory]{Items: items, Page: clampPage(page), PageSize: clampPageSize(size), Total: total})
}

func (h *AdminHandler) AdminCreateJobCategory(c *gin.Context) {
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
		logger.WithCtx(c).Error("catalog.admin.job_categories.create failed", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminPatchJobCategory(c *gin.Context) {
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
		logger.WithCtx(c).Error("catalog.admin.job_categories.patch failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminDeleteJobCategory(c *gin.Context) {
	if err := h.svc.repo.AdminDeleteJobCategory(c.Param("id")); err != nil {
		logger.WithCtx(c).Error("catalog.admin.job_categories.delete failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminListJobRoles(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("pageSize"), 20)
	q := c.Query("q")
	category := c.Query("category")
	items, total, err := h.svc.repo.AdminListJobRoles(page, size, q, category)
	if err != nil {
		logger.WithCtx(c).Error("catalog.admin.job_roles.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, AdminPageResp[JobRole]{Items: items, Page: clampPage(page), PageSize: clampPageSize(size), Total: total})
}

func (h *AdminHandler) AdminCreateJobRole(c *gin.Context) {
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
		logger.WithCtx(c).Error("catalog.admin.job_roles.create failed", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminPatchJobRole(c *gin.Context) {
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
		logger.WithCtx(c).Error("catalog.admin.job_roles.patch failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminDeleteJobRole(c *gin.Context) {
	if err := h.svc.repo.AdminDeleteJobRole(c.Param("id")); err != nil {
		logger.WithCtx(c).Error("catalog.admin.job_roles.delete failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminListContentPresets(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("pageSize"), 20)
	q := c.Query("q")
	role := c.Query("role")
	lang := c.Query("language")
	items, total, err := h.svc.repo.AdminListContentPresets(page, size, q, role, lang)
	if err != nil {
		logger.WithCtx(c).Error("catalog.admin.content_presets.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, AdminPageResp[ContentPreset]{Items: items, Page: clampPage(page), PageSize: clampPageSize(size), Total: total})
}

func (h *AdminHandler) AdminCreateContentPreset(c *gin.Context) {
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
		logger.WithCtx(c).Error("catalog.admin.content_presets.create failed", zap.Error(err))
		if strings.Contains(err.Error(), "invalid_json") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminPatchContentPreset(c *gin.Context) {
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
		logger.WithCtx(c).Error("catalog.admin.content_presets.patch failed", zap.Error(err))
		if strings.Contains(err.Error(), "invalid_json") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminDeleteContentPreset(c *gin.Context) {
	if err := h.svc.repo.AdminDeleteContentPreset(c.Param("id")); err != nil {
		logger.WithCtx(c).Error("catalog.admin.content_presets.delete failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminListTemplateVariants(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("pageSize"), 20)
	q := c.Query("q")
	role := c.Query("role")
	category := c.Query("category")
	template := c.Query("template")
	items, total, err := h.svc.repo.AdminListTemplateVariants(page, size, q, role, category, template)
	if err != nil {
		logger.WithCtx(c).Error("catalog.admin.template_variants.list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, AdminPageResp[TemplateVariant]{Items: items, Page: clampPage(page), PageSize: clampPageSize(size), Total: total})
}

func (h *AdminHandler) AdminCreateTemplateVariant(c *gin.Context) {
	var body struct {
		ExternalID               string `json:"externalId"`
		Name                     string `json:"name"`
		LayoutTemplateExternalID string `json:"layoutTemplateExternalId"`
		PresetExternalID         string `json:"presetExternalId"`
		RoleExternalID           string `json:"roleExternalId"`
		Tags                     string `json:"tags"`
		UsageCount               *int   `json:"usageCount"`
		IsPremium                *bool  `json:"isPremium"`
		IsActive                 *bool  `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.ExternalID) == "" || strings.TrimSpace(body.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	m := TemplateVariant{
		ExternalID:               strings.TrimSpace(body.ExternalID),
		Name:                     strings.TrimSpace(body.Name),
		LayoutTemplateExternalID: strings.TrimSpace(body.LayoutTemplateExternalID),
		PresetExternalID:         strings.TrimSpace(body.PresetExternalID),
		RoleExternalID:           strings.TrimSpace(body.RoleExternalID),
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
		logger.WithCtx(c).Error("catalog.admin.template_variants.create failed", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminPatchTemplateVariant(c *gin.Context) {
	var body struct {
		Name                     *string `json:"name"`
		LayoutTemplateExternalID *string `json:"layoutTemplateExternalId"`
		PresetExternalID         *string `json:"presetExternalId"`
		RoleExternalID           *string `json:"roleExternalId"`
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
	if body.PresetExternalID != nil {
		patch["preset_external_id"] = strings.TrimSpace(*body.PresetExternalID)
	}
	if body.RoleExternalID != nil {
		patch["role_external_id"] = strings.TrimSpace(*body.RoleExternalID)
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
	if err := h.svc.repo.AdminPatchTemplateVariant(c.Param("id"), patch); err != nil {
		logger.WithCtx(c).Error("catalog.admin.template_variants.patch failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminDeleteTemplateVariant(c *gin.Context) {
	if err := h.svc.repo.AdminDeleteTemplateVariant(c.Param("id")); err != nil {
		logger.WithCtx(c).Error("catalog.admin.template_variants.delete failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminGenerateTemplateVariants(c *gin.Context) {
	var req GenerateVariantsReq
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.RoleID) == "" || strings.TrimSpace(req.PresetID) == "" || len(req.LayoutTemplateIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	mode := GenerateModeSkip
	if strings.TrimSpace(req.Mode) == string(GenerateModeUpdate) {
		mode = GenerateModeUpdate
	}
	roleID := strings.TrimSpace(req.RoleID)
	presetID := strings.TrimSpace(req.PresetID)

	roleOk, err := h.svc.repo.existsActive(&JobRole{}, roleID)
	if err != nil || !roleOk {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}
	presetOk, err := h.svc.repo.existsActive(&ContentPreset{}, presetID)
	if err != nil || !presetOk {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid preset"})
		return
	}
	roleName, _ := h.svc.repo.getRoleName(roleID)
	presetName, _ := h.svc.repo.getPresetName(presetID)

	result := GenerateVariantsResult{}
	err = h.svc.repo.db.Transaction(func(tx *gorm.DB) error {
		for _, tplID := range req.LayoutTemplateIDs {
			tplID = strings.TrimSpace(tplID)
			item := struct {
				LayoutTemplateID string `json:"layoutTemplateId"`
				ExternalID       string `json:"externalId"`
				Action           string `json:"action"`
				Error            string `json:"error,omitempty"`
			}{LayoutTemplateID: tplID}

			if tplID == "" {
				item.Action = "failed"
				item.Error = "empty template"
				result.Failed++
				result.Items = append(result.Items, item)
				continue
			}

			tplOk, err := h.svc.repo.exists(&Template{}, tplID)
			if err != nil || !tplOk {
				item.Action = "failed"
				item.Error = "invalid template"
				result.Failed++
				result.Items = append(result.Items, item)
				continue
			}
			tplName, _ := h.svc.repo.getTemplateName(tplID)
			ext := stableExternalID("variant", roleID, presetID, tplID)
			item.ExternalID = ext
			name := strings.TrimSpace(req.NamePrefix)
			if name == "" {
				base := strings.TrimSpace(roleName)
				if base == "" {
					base = roleID
				}
				if tplName != "" {
					name = base + " - " + tplName
				} else {
					name = base + " - " + tplID
				}
			} else if tplName != "" {
				name = name + " - " + tplName
			} else {
				name = name + " - " + tplID
			}

			v := TemplateVariant{
				ExternalID:               ext,
				Name:                     name,
				LayoutTemplateExternalID: tplID,
				PresetExternalID:         presetID,
				RoleExternalID:           roleID,
				Tags:                     strings.TrimSpace(req.Tags),
				IsPremium:                false,
				IsActive:                 true,
			}
			if presetName != "" && v.Tags == "" {
				v.Tags = presetName
			}
			if req.IsPremium != nil {
				v.IsPremium = *req.IsPremium
			}
			if req.IsActive != nil {
				v.IsActive = *req.IsActive
			}

			action, err := h.svc.repo.upsertVariant(tx, &v, mode)
			if err != nil {
				item.Action = "failed"
				item.Error = err.Error()
				result.Failed++
				result.Items = append(result.Items, item)
				continue
			}
			item.Action = action
			switch action {
			case "created":
				result.Created++
			case "updated":
				result.Updated++
			default:
				result.Skipped++
			}
			result.Items = append(result.Items, item)
		}
		return nil
	})
	if err != nil {
		logger.WithCtx(c).Error("catalog.admin.template_variants.generate failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "result": result})
}
