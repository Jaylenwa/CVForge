package catalog

import (
	"errors"
	"strings"

	"openresume/internal/infra/database"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func DefaultRepo() *Repo {
	return &Repo{db: database.DB}
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) ListJobCategories() ([]JobCategory, error) {
	var list []JobCategory
	err := r.db.Where("is_active = ?", true).Order("order_num asc").Order("id asc").Find(&list).Error
	return list, err
}

func (r *Repo) ListJobRoles(categoryExternalID string, q string) ([]JobRole, error) {
	var list []JobRole
	db := r.db.Where("is_active = ?", true)
	if categoryExternalID != "" {
		db = db.Where("category_external_id = ?", categoryExternalID)
	}
	if q != "" {
		qq := "%" + strings.TrimSpace(q) + "%"
		db = db.Where("name LIKE ? OR tags LIKE ?", qq, qq)
	}
	err := db.Order("order_num asc").Order("id asc").Find(&list).Error
	return list, err
}

func (r *Repo) ListRoleExternalIDsByCategory(categoryExternalID string) ([]string, error) {
	if categoryExternalID == "" {
		return nil, nil
	}
	var roles []JobRole
	err := r.db.Select("external_id").Where("is_active = ?", true).Where("category_external_id = ?", categoryExternalID).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(roles))
	for _, rr := range roles {
		out = append(out, rr.ExternalID)
	}
	return out, nil
}

func (r *Repo) ListTemplateVariants(roleExternalID string, categoryExternalID string, q string) ([]TemplateVariant, error) {
	var list []TemplateVariant
	db := r.db.Where("is_active = ?", true)
	if roleExternalID != "" {
		db = db.Where("role_external_id = ?", roleExternalID)
	} else if categoryExternalID != "" {
		roleIDs, err := r.ListRoleExternalIDsByCategory(categoryExternalID)
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

func (r *Repo) GetContentPresetByExternal(id string) (ContentPreset, error) {
	var p ContentPreset
	err := r.db.Where("external_id = ?", id).Where("is_active = ?", true).First(&p).Error
	return p, err
}

func (r *Repo) UpsertJobCategory(db *gorm.DB, c *JobCategory) error {
	if c == nil || c.ExternalID == "" {
		return errors.New("invalid job category")
	}
	var existing JobCategory
	err := db.Where("external_id = ?", c.ExternalID).First(&existing).Error
	if err == nil {
		existing.Name = c.Name
		existing.ParentExternalID = c.ParentExternalID
		existing.OrderNum = c.OrderNum
		existing.IsActive = c.IsActive
		return db.Save(&existing).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return db.Create(c).Error
	}
	return err
}

func (r *Repo) UpsertJobRole(db *gorm.DB, rr *JobRole) error {
	if rr == nil || rr.ExternalID == "" {
		return errors.New("invalid job role")
	}
	var existing JobRole
	err := db.Where("external_id = ?", rr.ExternalID).First(&existing).Error
	if err == nil {
		existing.CategoryExternalID = rr.CategoryExternalID
		existing.Name = rr.Name
		existing.Tags = rr.Tags
		existing.OrderNum = rr.OrderNum
		existing.IsActive = rr.IsActive
		return db.Save(&existing).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return db.Create(rr).Error
	}
	return err
}

func (r *Repo) UpsertContentPreset(db *gorm.DB, p *ContentPreset) error {
	if p == nil || p.ExternalID == "" {
		return errors.New("invalid content preset")
	}
	var existing ContentPreset
	err := db.Where("external_id = ?", p.ExternalID).First(&existing).Error
	if err == nil {
		existing.Name = p.Name
		existing.Language = p.Language
		existing.RoleExternalID = p.RoleExternalID
		existing.Tags = p.Tags
		existing.DataJSON = p.DataJSON
		existing.IsActive = p.IsActive
		return db.Save(&existing).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return db.Create(p).Error
	}
	return err
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
