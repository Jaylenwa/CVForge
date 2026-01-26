package library

import (
	"errors"
	"strconv"
	"strings"

	"openresume/internal/infra/database"
	"openresume/internal/module/taxonomy"

	"gorm.io/gorm"
)

type Repo struct {
	db      *gorm.DB
	taxRepo *taxonomy.Repo
}

func DefaultRepo() *Repo {
	return &Repo{db: database.DB, taxRepo: taxonomy.DefaultRepo()}
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db, taxRepo: taxonomy.NewRepo(db)}
}

func (r *Repo) ListTemplateVariants(roleID uint, categoryID uint, q string) ([]TemplateVariant, error) {
	var list []TemplateVariant
	db := r.db.Where("is_active = ?", true)
	if roleID != 0 {
		db = db.Where("role_id = ?", roleID)
	} else if categoryID != 0 {
		roleIDs, err := r.taxRepo.ListRoleIDsByCategory(categoryID)
		if err != nil {
			return nil, err
		}
		if len(roleIDs) == 0 {
			return []TemplateVariant{}, nil
		}
		db = db.Where("role_id IN ?", roleIDs)
	}
	if q != "" {
		qq := "%" + strings.TrimSpace(q) + "%"
		db = db.Where("name LIKE ? OR tags LIKE ?", qq, qq)
	}
	err := db.Order("usage_count desc").Order("id asc").Find(&list).Error
	return list, err
}

func (r *Repo) GetTemplateVariantByID(id uint) (TemplateVariant, error) {
	var v TemplateVariant
	err := r.db.Where("id = ?", id).Where("is_active = ?", true).First(&v).Error
	return v, err
}

func (r *Repo) UpsertTemplateVariant(db *gorm.DB, v *TemplateVariant) error {
	if v == nil || v.RoleID == 0 || v.PresetID == 0 || strings.TrimSpace(v.LayoutTemplateExternalID) == "" {
		return errors.New("invalid template variant")
	}
	var existing TemplateVariant
	err := db.Where("role_id = ?", v.RoleID).
		Where("preset_id = ?", v.PresetID).
		Where("layout_template_external_id = ?", v.LayoutTemplateExternalID).
		First(&existing).Error
	if err == nil {
		v.ID = existing.ID
		existing.Name = v.Name
		existing.LayoutTemplateExternalID = v.LayoutTemplateExternalID
		existing.PresetID = v.PresetID
		existing.RoleID = v.RoleID
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

func (r *Repo) IncrementTemplateVariantUsage(id uint) error {
	if id == 0 {
		return nil
	}
	return r.db.Model(&TemplateVariant{}).Where("id = ?", id).UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).Error
}

func (r *Repo) AdminListTemplateVariants(page, size int, q, role, category, template string) ([]TemplateVariant, int64, error) {
	page = clampPage(page)
	size = clampPageSize(size)
	var list []TemplateVariant
	db := r.db.Model(&TemplateVariant{})
	if q = strings.TrimSpace(q); q != "" {
		qq := "%" + q + "%"
		db = db.Where("name LIKE ? OR tags LIKE ?", qq, qq)
		if id, err := strconv.ParseUint(q, 10, 64); err == nil {
			db = db.Or("id = ?", uint(id))
		}
	}
	if role = strings.TrimSpace(role); role != "" {
		if roleID, err := strconv.ParseUint(role, 10, 64); err == nil {
			db = db.Where("role_id = ?", uint(roleID))
		}
	} else if category = strings.TrimSpace(category); category != "" {
		catID, err := strconv.ParseUint(category, 10, 64)
		if err != nil {
			return []TemplateVariant{}, 0, nil
		}
		roleIDs, err := r.taxRepo.ListRoleIDsByCategory(uint(catID))
		if err != nil {
			return nil, 0, err
		}
		if len(roleIDs) == 0 {
			return []TemplateVariant{}, 0, nil
		}
		db = db.Where("role_id IN ?", roleIDs)
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
	if v == nil || strings.TrimSpace(v.Name) == "" || v.RoleID == 0 || v.PresetID == 0 || strings.TrimSpace(v.LayoutTemplateExternalID) == "" {
		return errors.New("invalid")
	}
	return r.db.Create(v).Error
}

func (r *Repo) AdminPatchTemplateVariant(id uint, patch map[string]any) error {
	return r.db.Model(&TemplateVariant{}).Where("id = ?", id).Updates(patch).Error
}

func (r *Repo) AdminDeleteTemplateVariant(id uint) error {
	return r.db.Where("id = ?", id).Delete(&TemplateVariant{}).Error
}

func (r *Repo) getTemplateName(externalID string) (string, error) {
	var t struct {
		Name string
	}
	err := r.db.Table("template").Select("name").Where("external_id = ?", externalID).Where("deleted_at IS NULL").Scan(&t).Error
	return t.Name, err
}

func (r *Repo) getRoleName(id uint) (string, error) {
	var t struct {
		Name string
	}
	err := r.db.Table("job_role").Select("name").Where("id = ?", id).Where("deleted_at IS NULL").Scan(&t).Error
	return t.Name, err
}

func (r *Repo) getPresetName(id uint) (string, error) {
	var t struct {
		Name string
	}
	err := r.db.Table("content_preset").Select("name").Where("id = ?", id).Where("deleted_at IS NULL").Scan(&t).Error
	return t.Name, err
}

func (r *Repo) findVariantByCombo(roleID, presetID uint, templateExternalID string) (TemplateVariant, error) {
	var v TemplateVariant
	err := r.db.Where("role_id = ?", roleID).
		Where("preset_id = ?", presetID).
		Where("layout_template_external_id = ?", templateExternalID).
		First(&v).Error
	return v, err
}

func (r *Repo) upsertVariant(tx *gorm.DB, v *TemplateVariant, mode GenerateMode) (string, error) {
	existing, err := r.findVariantByCombo(v.RoleID, v.PresetID, v.LayoutTemplateExternalID)
	if err == nil {
		if mode != GenerateModeUpdate {
			return "skipped", nil
		}
		existing.Name = v.Name
		existing.LayoutTemplateExternalID = v.LayoutTemplateExternalID
		existing.PresetID = v.PresetID
		existing.RoleID = v.RoleID
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
