package library

import (
	"errors"
	"strings"

	"openresume/internal/infra/database"
	"openresume/internal/module/taxonomy"

	"gorm.io/gorm"
)

type Repo struct {
	db     *gorm.DB
	taxRepo *taxonomy.Repo
}

func DefaultRepo() *Repo {
	return &Repo{db: database.DB, taxRepo: taxonomy.DefaultRepo()}
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db, taxRepo: taxonomy.NewRepo(db)}
}

func (r *Repo) ListTemplateVariants(roleExternalID string, categoryExternalID string, q string) ([]TemplateVariant, error) {
	var list []TemplateVariant
	db := r.db.Where("is_active = ?", true)
	if roleExternalID != "" {
		db = db.Where("role_external_id = ?", roleExternalID)
	} else if categoryExternalID != "" {
		roleIDs, err := r.taxRepo.ListRoleExternalIDsByCategory(categoryExternalID)
		if err != nil {
			return nil, err
		}
		if len(roleIDs) == 0 {
			return []TemplateVariant{}, nil
		}
		db = db.Where("role_external_id IN ?", roleIDs)
	}
	if q != "" {
		qq := "%" + strings.TrimSpace(q) + "%"
		db = db.Where("name LIKE ? OR tags LIKE ?", qq, qq)
	}
	err := db.Order("usage_count desc").Order("id asc").Find(&list).Error
	return list, err
}

func (r *Repo) GetTemplateVariantByExternal(externalID string) (TemplateVariant, error) {
	var v TemplateVariant
	err := r.db.Where("external_id = ?", externalID).Where("is_active = ?", true).First(&v).Error
	return v, err
}

func (r *Repo) UpsertTemplateVariant(db *gorm.DB, v *TemplateVariant) error {
	if v == nil || v.ExternalID == "" {
		return errors.New("invalid template variant")
	}
	var existing TemplateVariant
	err := db.Where("external_id = ?", v.ExternalID).First(&existing).Error
	if err == nil {
		existing.Name = v.Name
		existing.LayoutTemplateExternalID = v.LayoutTemplateExternalID
		existing.PresetExternalID = v.PresetExternalID
		existing.RoleExternalID = v.RoleExternalID
		existing.Tags = v.Tags
		existing.UsageCount = v.UsageCount
		existing.IsPremium = v.IsPremium
		existing.IsActive = v.IsActive
		return db.Save(&existing).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return db.Create(v).Error
	}
	return err
}

func (r *Repo) IncrementTemplateVariantUsage(externalID string) error {
	if externalID == "" {
		return nil
	}
	return r.db.Model(&TemplateVariant{}).Where("external_id = ?", externalID).UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).Error
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
		roleIDs, err := r.taxRepo.ListRoleExternalIDsByCategory(category)
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

