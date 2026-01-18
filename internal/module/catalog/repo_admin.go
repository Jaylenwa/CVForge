package catalog

import (
	"encoding/json"
	"errors"
	"strings"

	"gorm.io/gorm"
)

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

func (r *Repo) AdminListJobCategories(page, size int, q string, parent string) ([]JobCategory, int64, error) {
	page = clampPage(page)
	size = clampPageSize(size)
	var list []JobCategory
	db := r.db.Model(&JobCategory{})
	if q = strings.TrimSpace(q); q != "" {
		qq := "%" + q + "%"
		db = db.Where("name LIKE ? OR external_id LIKE ?", qq, qq)
	}
	if parent = strings.TrimSpace(parent); parent != "" {
		db = db.Where("parent_external_id = ?", parent)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("order_num asc").Order("id asc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *Repo) AdminGetJobCategory(externalID string) (JobCategory, error) {
	var c JobCategory
	err := r.db.Where("external_id = ?", externalID).First(&c).Error
	return c, err
}

func (r *Repo) AdminCreateJobCategory(c *JobCategory) error {
	if c == nil || strings.TrimSpace(c.ExternalID) == "" || strings.TrimSpace(c.Name) == "" {
		return errors.New("invalid")
	}
	return r.db.Create(c).Error
}

func (r *Repo) AdminPatchJobCategory(externalID string, patch map[string]any) error {
	return r.db.Model(&JobCategory{}).Where("external_id = ?", externalID).Updates(patch).Error
}

func (r *Repo) AdminDeleteJobCategory(externalID string) error {
	return r.db.Where("external_id = ?", externalID).Delete(&JobCategory{}).Error
}

func (r *Repo) AdminListJobRoles(page, size int, q, category string) ([]JobRole, int64, error) {
	page = clampPage(page)
	size = clampPageSize(size)
	var list []JobRole
	db := r.db.Model(&JobRole{})
	if q = strings.TrimSpace(q); q != "" {
		qq := "%" + q + "%"
		db = db.Where("name LIKE ? OR tags LIKE ? OR external_id LIKE ?", qq, qq, qq)
	}
	if category = strings.TrimSpace(category); category != "" {
		db = db.Where("category_external_id = ?", category)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("order_num asc").Order("id asc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *Repo) AdminCreateJobRole(rr *JobRole) error {
	if rr == nil || strings.TrimSpace(rr.ExternalID) == "" || strings.TrimSpace(rr.Name) == "" || strings.TrimSpace(rr.CategoryExternalID) == "" {
		return errors.New("invalid")
	}
	return r.db.Create(rr).Error
}

func (r *Repo) AdminPatchJobRole(externalID string, patch map[string]any) error {
	return r.db.Model(&JobRole{}).Where("external_id = ?", externalID).Updates(patch).Error
}

func (r *Repo) AdminDeleteJobRole(externalID string) error {
	return r.db.Where("external_id = ?", externalID).Delete(&JobRole{}).Error
}

func (r *Repo) AdminListContentPresets(page, size int, q, role, language string) ([]ContentPreset, int64, error) {
	page = clampPage(page)
	size = clampPageSize(size)
	var list []ContentPreset
	db := r.db.Model(&ContentPreset{})
	if q = strings.TrimSpace(q); q != "" {
		qq := "%" + q + "%"
		db = db.Where("name LIKE ? OR tags LIKE ? OR external_id LIKE ?", qq, qq, qq)
	}
	if role = strings.TrimSpace(role); role != "" {
		db = db.Where("role_external_id = ?", role)
	}
	if language = strings.TrimSpace(language); language != "" {
		db = db.Where("language = ?", language)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *Repo) AdminCreateContentPreset(p *ContentPreset) error {
	if p == nil || strings.TrimSpace(p.ExternalID) == "" || strings.TrimSpace(p.Name) == "" {
		return errors.New("invalid")
	}
	if p.DataJSON != "" && !json.Valid([]byte(p.DataJSON)) {
		return errors.New("invalid_json")
	}
	return r.db.Create(p).Error
}

func (r *Repo) AdminPatchContentPreset(externalID string, patch map[string]any) error {
	if v, ok := patch["data_json"]; ok {
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" && !json.Valid([]byte(s)) {
			return errors.New("invalid_json")
		}
	}
	return r.db.Model(&ContentPreset{}).Where("external_id = ?", externalID).Updates(patch).Error
}

func (r *Repo) AdminDeleteContentPreset(externalID string) error {
	return r.db.Where("external_id = ?", externalID).Delete(&ContentPreset{}).Error
}

func (r *Repo) AdminListTemplateVariants(page, size int, q, role, category, template string) ([]TemplateVariant, int64, error) {
	page = clampPage(page)
	size = clampPageSize(size)
	var list []TemplateVariant
	db := r.db.Model(&TemplateVariant{})
	if q = strings.TrimSpace(q); q != "" {
		qq := "%" + q + "%"
		db = db.Where("name LIKE ? OR tags LIKE ? OR external_id LIKE ?", qq, qq, qq)
	}
	if role = strings.TrimSpace(role); role != "" {
		db = db.Where("role_external_id = ?", role)
	} else if category = strings.TrimSpace(category); category != "" {
		roleIDs, err := r.ListRoleExternalIDsByCategory(category)
		if err != nil {
			return nil, 0, err
		}
		if len(roleIDs) == 0 {
			return []TemplateVariant{}, 0, nil
		}
		db = db.Where("role_external_id IN ?", roleIDs)
	}
	if template = strings.TrimSpace(template); template != "" {
		db = db.Where("layout_template_external_id = ?", template)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *Repo) AdminCreateTemplateVariant(v *TemplateVariant) error {
	if v == nil || strings.TrimSpace(v.ExternalID) == "" || strings.TrimSpace(v.Name) == "" {
		return errors.New("invalid")
	}
	return r.db.Create(v).Error
}

func (r *Repo) AdminPatchTemplateVariant(externalID string, patch map[string]any) error {
	return r.db.Model(&TemplateVariant{}).Where("external_id = ?", externalID).Updates(patch).Error
}

func (r *Repo) AdminDeleteTemplateVariant(externalID string) error {
	return r.db.Where("external_id = ?", externalID).Delete(&TemplateVariant{}).Error
}

func (r *Repo) existsActive(model any, externalID string) (bool, error) {
	if externalID == "" {
		return false, nil
	}
	var cnt int64
	err := r.db.Model(model).Where("external_id = ?", externalID).Where("is_active = ?", true).Count(&cnt).Error
	return cnt > 0, err
}

func (r *Repo) exists(model any, externalID string) (bool, error) {
	if externalID == "" {
		return false, nil
	}
	var cnt int64
	err := r.db.Model(model).Where("external_id = ?", externalID).Count(&cnt).Error
	return cnt > 0, err
}

func (r *Repo) getTemplateName(externalID string) (string, error) {
	var t struct {
		Name string
	}
	err := r.db.Table("template").Select("name").Where("external_id = ?", externalID).Where("deleted_at IS NULL").Scan(&t).Error
	return t.Name, err
}

func (r *Repo) getRoleName(externalID string) (string, error) {
	var t struct {
		Name string
	}
	err := r.db.Table("job_role").Select("name").Where("external_id = ?", externalID).Where("deleted_at IS NULL").Scan(&t).Error
	return t.Name, err
}

func (r *Repo) getPresetName(externalID string) (string, error) {
	var t struct {
		Name string
	}
	err := r.db.Table("content_preset").Select("name").Where("external_id = ?", externalID).Where("deleted_at IS NULL").Scan(&t).Error
	return t.Name, err
}

func (r *Repo) findVariantByExternal(externalID string) (TemplateVariant, error) {
	var v TemplateVariant
	err := r.db.Where("external_id = ?", externalID).First(&v).Error
	return v, err
}

func (r *Repo) upsertVariant(tx *gorm.DB, v *TemplateVariant, mode GenerateMode) (string, error) {
	existing, err := r.findVariantByExternal(v.ExternalID)
	if err == nil {
		if mode != GenerateModeUpdate {
			return "skipped", nil
		}
		existing.Name = v.Name
		existing.LayoutTemplateExternalID = v.LayoutTemplateExternalID
		existing.PresetExternalID = v.PresetExternalID
		existing.RoleExternalID = v.RoleExternalID
		existing.Tags = v.Tags
		existing.IsPremium = v.IsPremium
		existing.IsActive = v.IsActive
		return "updated", tx.Save(&existing).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "created", tx.Create(v).Error
	}
	return "", err
}
