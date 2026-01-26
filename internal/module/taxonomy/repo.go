package taxonomy

import (
	"errors"
	"strconv"
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

func (r *Repo) ListJobRoles(categoryID uint, q string) ([]JobRole, error) {
	var list []JobRole
	db := r.db.Where("is_active = ?", true)
	if categoryID != 0 {
		db = db.Where("category_id = ?", categoryID)
	}
	if q != "" {
		qq := "%" + strings.TrimSpace(q) + "%"
		db = db.Where("name LIKE ? OR tags LIKE ?", qq, qq)
	}
	err := db.Order("order_num asc").Order("id asc").Find(&list).Error
	return list, err
}

func (r *Repo) ListRoleIDsByCategory(categoryID uint) ([]uint, error) {
	if categoryID == 0 {
		return nil, nil
	}
	var roles []JobRole
	err := r.db.Select("id").Where("is_active = ?", true).Where("category_id = ?", categoryID).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	out := make([]uint, 0, len(roles))
	for _, rr := range roles {
		out = append(out, rr.ID)
	}
	return out, nil
}

func (r *Repo) UpsertJobCategory(db *gorm.DB, c *JobCategory) error {
	if c == nil || strings.TrimSpace(c.Name) == "" {
		return errors.New("invalid job category")
	}
	var existing JobCategory
	err := db.Where("name = ?", strings.TrimSpace(c.Name)).Where("parent_id IS ?", c.ParentID).First(&existing).Error
	if err == nil {
		c.ID = existing.ID
		existing.Name = c.Name
		existing.ParentID = c.ParentID
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
	if rr == nil || strings.TrimSpace(rr.Name) == "" || rr.CategoryID == 0 {
		return errors.New("invalid job role")
	}
	var existing JobRole
	err := db.Where("category_id = ?", rr.CategoryID).Where("name = ?", strings.TrimSpace(rr.Name)).First(&existing).Error
	if err == nil {
		rr.ID = existing.ID
		existing.CategoryID = rr.CategoryID
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

func (r *Repo) AdminListJobCategories(page, size int, q string, parentID *uint) ([]JobCategory, int64, error) {
	page = clampPage(page)
	size = clampPageSize(size)
	var list []JobCategory
	db := r.db.Model(&JobCategory{})
	if q = strings.TrimSpace(q); q != "" {
		qq := "%" + q + "%"
		db = db.Where("name LIKE ?", qq)
		if id, err := strconv.ParseUint(q, 10, 64); err == nil {
			db = db.Or("id = ?", uint(id))
		}
	}
	if parentID != nil {
		db = db.Where("parent_id = ?", *parentID)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("order_num asc").Order("id asc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *Repo) AdminCreateJobCategory(c *JobCategory) error {
	if c == nil || strings.TrimSpace(c.Name) == "" {
		return errors.New("invalid")
	}
	return r.db.Create(c).Error
}

func (r *Repo) AdminPatchJobCategory(id uint, patch map[string]any) error {
	return r.db.Model(&JobCategory{}).Where("id = ?", id).Updates(patch).Error
}

func (r *Repo) AdminDeleteJobCategory(id uint) error {
	return r.db.Where("id = ?", id).Delete(&JobCategory{}).Error
}

func (r *Repo) AdminListJobRoles(page, size int, q string, categoryID *uint) ([]JobRole, int64, error) {
	page = clampPage(page)
	size = clampPageSize(size)
	var list []JobRole
	db := r.db.Model(&JobRole{})
	if q = strings.TrimSpace(q); q != "" {
		qq := "%" + q + "%"
		db = db.Where("name LIKE ? OR tags LIKE ?", qq, qq)
		if id, err := strconv.ParseUint(q, 10, 64); err == nil {
			db = db.Or("id = ?", uint(id))
		}
	}
	if categoryID != nil {
		db = db.Where("category_id = ?", *categoryID)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("order_num asc").Order("id asc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *Repo) AdminCreateJobRole(rr *JobRole) error {
	if rr == nil || rr.CategoryID == 0 || strings.TrimSpace(rr.Name) == "" {
		return errors.New("invalid")
	}
	return r.db.Create(rr).Error
}

func (r *Repo) AdminPatchJobRole(id uint, patch map[string]any) error {
	return r.db.Model(&JobRole{}).Where("id = ?", id).Updates(patch).Error
}

func (r *Repo) AdminDeleteJobRole(id uint) error {
	return r.db.Where("id = ?", id).Delete(&JobRole{}).Error
}
