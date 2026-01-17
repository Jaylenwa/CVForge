package catalog

import (
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

func (r *Repo) GetContentPresetByExternal(id string) (ContentPreset, error) {
	var p ContentPreset
	err := r.db.Where("external_id = ?", id).Where("is_active = ?", true).First(&p).Error
	return p, err
}

