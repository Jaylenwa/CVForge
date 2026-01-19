package taxonomy

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

