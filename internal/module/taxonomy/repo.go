package taxonomy

import (
	"errors"
	"strconv"
	"strings"

	"openresume/internal/infra/database"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *Repo) ListJobCategories(language string) ([]JobCategoryView, error) {
	var list []JobCategoryView
	language = normalizeLanguage(language)
	err := r.db.Model(&JobCategory{}).
		Select("job_category.*, COALESCE(job_category_i18n_req.name, job_category_i18n_zh.name) as name").
		Joins("LEFT JOIN job_category_i18n job_category_i18n_req ON job_category_i18n_req.job_category_id = job_category.id AND job_category_i18n_req.language = ? AND job_category_i18n_req.deleted_at IS NULL", language).
		Joins("LEFT JOIN job_category_i18n job_category_i18n_zh ON job_category_i18n_zh.job_category_id = job_category.id AND job_category_i18n_zh.language = 'zh' AND job_category_i18n_zh.deleted_at IS NULL").
		Where("job_category.is_active = ?", true).
		Order("job_category.order_num asc").
		Order("job_category.id asc").
		Find(&list).Error
	return list, err
}

func (r *Repo) ListJobRoles(language string, categoryID uint, q string) ([]JobRoleView, error) {
	var list []JobRoleView
	language = normalizeLanguage(language)
	db := r.db.Model(&JobRole{}).
		Select("job_role.*, COALESCE(job_role_i18n_req.name, job_role_i18n_zh.name) as name").
		Joins("LEFT JOIN job_role_i18n job_role_i18n_req ON job_role_i18n_req.job_role_id = job_role.id AND job_role_i18n_req.language = ? AND job_role_i18n_req.deleted_at IS NULL", language).
		Joins("LEFT JOIN job_role_i18n job_role_i18n_zh ON job_role_i18n_zh.job_role_id = job_role.id AND job_role_i18n_zh.language = 'zh' AND job_role_i18n_zh.deleted_at IS NULL").
		Where("job_role.is_active = ?", true)
	if categoryID != 0 {
		db = db.Where("job_role.category_id = ?", categoryID)
	}
	if q != "" {
		qq := "%" + strings.TrimSpace(q) + "%"
		db = db.Where("COALESCE(job_role_i18n_req.name, job_role_i18n_zh.name) LIKE ?", qq)
	}
	err := db.Order("job_role.order_num asc").Order("job_role.id asc").Find(&list).Error
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

func (r *Repo) ListJobCategoryI18n(categoryIDs []uint) ([]JobCategoryI18n, error) {
	if len(categoryIDs) == 0 {
		return nil, nil
	}
	var list []JobCategoryI18n
	err := r.db.Where("job_category_id IN ?", categoryIDs).Find(&list).Error
	return list, err
}

func (r *Repo) ListJobRoleI18n(roleIDs []uint) ([]JobRoleI18n, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}
	var list []JobRoleI18n
	err := r.db.Where("job_role_id IN ?", roleIDs).Find(&list).Error
	return list, err
}

func upsertJobCategoryI18n(tx *gorm.DB, categoryID uint, language, name string) error {
	language = strings.ToLower(strings.TrimSpace(language))
	if language != "zh" && language != "en" {
		return nil
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	m := JobCategoryI18n{
		JobCategoryID: categoryID,
		Language:      language,
		Name:          name,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "job_category_id"}, {Name: "language"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&m).Error
}

func upsertJobRoleI18n(tx *gorm.DB, roleID uint, language, name string) error {
	language = strings.ToLower(strings.TrimSpace(language))
	if language != "zh" && language != "en" {
		return nil
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	m := JobRoleI18n{
		JobRoleID: roleID,
		Language:  language,
		Name:      name,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "job_role_id"}, {Name: "language"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&m).Error
}

func (r *Repo) AdminCreateJobCategoryWithNames(c *JobCategory, names map[string]string) error {
	if c == nil {
		return errors.New("invalid")
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(c).Error; err != nil {
			return err
		}
		for lang, name := range names {
			if err := upsertJobCategoryI18n(tx, c.ID, lang, name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repo) AdminPatchJobCategoryWithNames(id uint, patch map[string]any, names map[string]string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if len(patch) > 0 {
			if err := tx.Model(&JobCategory{}).Where("id = ?", id).Updates(patch).Error; err != nil {
				return err
			}
		}
		for lang, name := range names {
			if err := upsertJobCategoryI18n(tx, id, lang, name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repo) AdminCreateJobRoleWithNames(rr *JobRole, names map[string]string) error {
	if rr == nil || rr.CategoryID == 0 {
		return errors.New("invalid")
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(rr).Error; err != nil {
			return err
		}
		for lang, name := range names {
			if err := upsertJobRoleI18n(tx, rr.ID, lang, name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repo) AdminPatchJobRoleWithNames(id uint, patch map[string]any, names map[string]string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if len(patch) > 0 {
			if err := tx.Model(&JobRole{}).Where("id = ?", id).Updates(patch).Error; err != nil {
				return err
			}
		}
		for lang, name := range names {
			if err := upsertJobRoleI18n(tx, id, lang, name); err != nil {
				return err
			}
		}
		return nil
	})
}

func pickAnyName(names map[string]string) string {
	if names == nil {
		return ""
	}
	if v := strings.TrimSpace(names["zh"]); v != "" {
		return v
	}
	if v := strings.TrimSpace(names["en"]); v != "" {
		return v
	}
	for _, v := range names {
		if v = strings.TrimSpace(v); v != "" {
			return v
		}
	}
	return ""
}

func (r *Repo) UpsertJobCategoryWithNames(db *gorm.DB, c *JobCategory, names map[string]string) error {
	if c == nil {
		return errors.New("invalid job category")
	}
	if c.ExternalID == nil {
		return errors.New("invalid job category external_id")
	}
	externalID := strings.TrimSpace(*c.ExternalID)
	if externalID == "" {
		return errors.New("invalid job category external_id")
	}
	c.ExternalID = &externalID
	keyName := pickAnyName(names)
	if keyName == "" {
		return errors.New("invalid job category name")
	}
	var existing JobCategory
	err := db.Model(&JobCategory{}).
		Where("external_id = ?", externalID).
		First(&existing).Error
	if err == nil {
		c.ID = existing.ID
		existing.ExternalID = &externalID
		existing.ParentID = c.ParentID
		existing.OrderNum = c.OrderNum
		existing.IsActive = c.IsActive
		if err := db.Save(&existing).Error; err != nil {
			return err
		}
		for lang, name := range names {
			if err := upsertJobCategoryI18n(db, existing.ID, lang, name); err != nil {
				return err
			}
		}
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	var byName JobCategory
	q := db.Model(&JobCategory{}).
		Joins("JOIN job_category_i18n ON job_category_i18n.job_category_id = job_category.id AND job_category_i18n.language = 'zh' AND job_category_i18n.deleted_at IS NULL").
		Where("job_category_i18n.name = ?", keyName)
	if c.ParentID == nil {
		q = q.Where("job_category.parent_id IS NULL")
	} else {
		q = q.Where("job_category.parent_id = ?", *c.ParentID)
	}
	if err := q.First(&byName).Error; err == nil {
		if byName.ExternalID == nil || strings.TrimSpace(*byName.ExternalID) == "" {
			c.ID = byName.ID
			byName.ExternalID = &externalID
			byName.ParentID = c.ParentID
			byName.OrderNum = c.OrderNum
			byName.IsActive = c.IsActive
			if err := db.Save(&byName).Error; err != nil {
				return err
			}
			for lang, name := range names {
				if err := upsertJobCategoryI18n(db, byName.ID, lang, name); err != nil {
					return err
				}
			}
			return nil
		}
	}
	if err := db.Create(c).Error; err != nil {
		return err
	}
	if names == nil {
		names = map[string]string{}
	}
	if strings.TrimSpace(names["zh"]) == "" {
		names["zh"] = keyName
	}
	for lang, name := range names {
		if err := upsertJobCategoryI18n(db, c.ID, lang, name); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) UpsertJobRoleWithNames(db *gorm.DB, rr *JobRole, names map[string]string) error {
	if rr == nil || rr.CategoryID == 0 {
		return errors.New("invalid job role")
	}
	if rr.ExternalID == nil {
		return errors.New("invalid job role external_id")
	}
	externalID := strings.TrimSpace(*rr.ExternalID)
	if externalID == "" {
		return errors.New("invalid job role external_id")
	}
	rr.ExternalID = &externalID
	keyName := pickAnyName(names)
	if keyName == "" {
		return errors.New("invalid job role name")
	}
	var existing JobRole
	err := db.Model(&JobRole{}).
		Where("external_id = ?", externalID).
		First(&existing).Error
	if err == nil {
		rr.ID = existing.ID
		existing.ExternalID = &externalID
		existing.CategoryID = rr.CategoryID
		existing.OrderNum = rr.OrderNum
		existing.IsActive = rr.IsActive
		if err := db.Save(&existing).Error; err != nil {
			return err
		}
		for lang, name := range names {
			if err := upsertJobRoleI18n(db, existing.ID, lang, name); err != nil {
				return err
			}
		}
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	var byName JobRole
	q := db.Model(&JobRole{}).
		Joins("JOIN job_role_i18n ON job_role_i18n.job_role_id = job_role.id AND job_role_i18n.language = 'zh' AND job_role_i18n.deleted_at IS NULL").
		Where("job_role.category_id = ?", rr.CategoryID).
		Where("job_role_i18n.name = ?", keyName)
	if err := q.First(&byName).Error; err == nil {
		if byName.ExternalID == nil || strings.TrimSpace(*byName.ExternalID) == "" {
			rr.ID = byName.ID
			byName.ExternalID = &externalID
			byName.CategoryID = rr.CategoryID
			byName.OrderNum = rr.OrderNum
			byName.IsActive = rr.IsActive
			if err := db.Save(&byName).Error; err != nil {
				return err
			}
			for lang, name := range names {
				if err := upsertJobRoleI18n(db, byName.ID, lang, name); err != nil {
					return err
				}
			}
			return nil
		}
	}
	if err := db.Create(rr).Error; err != nil {
		return err
	}
	if names == nil {
		names = map[string]string{}
	}
	if strings.TrimSpace(names["zh"]) == "" {
		names["zh"] = keyName
	}
	for lang, name := range names {
		if err := upsertJobRoleI18n(db, rr.ID, lang, name); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) AdminListJobCategories(page, size int, q string, parentID *uint) ([]JobCategory, int64, error) {
	page = clampPage(page)
	size = clampPageSize(size)
	var list []JobCategory
	baseDB := r.db.Model(&JobCategory{})
	if q = strings.TrimSpace(q); q != "" {
		qq := "%" + q + "%"
		baseDB = baseDB.
			Joins("LEFT JOIN job_category_i18n ON job_category_i18n.job_category_id = job_category.id AND job_category_i18n.deleted_at IS NULL").
			Where("job_category_i18n.name LIKE ?", qq)
		if id, err := strconv.ParseUint(q, 10, 64); err == nil {
			baseDB = baseDB.Or("job_category.id = ?", uint(id))
		}
	}
	if parentID != nil {
		baseDB = baseDB.Where("job_category.parent_id = ?", *parentID)
	}
	var total int64
	countDB := baseDB
	if strings.TrimSpace(q) != "" {
		countDB = countDB.Distinct("job_category.id")
	}
	if err := countDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	listDB := baseDB
	if strings.TrimSpace(q) != "" {
		listDB = listDB.Distinct().Select("job_category.*")
	}
	err := listDB.Order("job_category.order_num asc").Order("job_category.id asc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *Repo) AdminCreateJobCategory(c *JobCategory) error {
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
	baseDB := r.db.Model(&JobRole{})
	if q = strings.TrimSpace(q); q != "" {
		qq := "%" + q + "%"
		baseDB = baseDB.
			Joins("LEFT JOIN job_role_i18n ON job_role_i18n.job_role_id = job_role.id AND job_role_i18n.deleted_at IS NULL").
			Where("job_role_i18n.name LIKE ?", qq)
		if id, err := strconv.ParseUint(q, 10, 64); err == nil {
			baseDB = baseDB.Or("job_role.id = ?", uint(id))
		}
	}
	if categoryID != nil {
		baseDB = baseDB.Where("job_role.category_id = ?", *categoryID)
	}
	var total int64
	countDB := baseDB
	if strings.TrimSpace(q) != "" {
		countDB = countDB.Distinct("job_role.id")
	}
	if err := countDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	listDB := baseDB
	if strings.TrimSpace(q) != "" {
		listDB = listDB.Distinct().Select("job_role.*")
	}
	err := listDB.Order("job_role.order_num asc").Order("job_role.id asc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *Repo) AdminCreateJobRole(rr *JobRole) error {
	if rr == nil || rr.CategoryID == 0 {
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
