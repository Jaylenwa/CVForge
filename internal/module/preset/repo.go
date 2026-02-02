package preset

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

func (r *Repo) GetByIDActive(id uint) (ContentPreset, error) {
	var p ContentPreset
	err := r.db.Where("id = ?", id).Where("is_active = ?", true).First(&p).Error
	return p, err
}

func (r *Repo) UpsertContentPreset(db *gorm.DB, p *ContentPreset) error {
	if p == nil || strings.TrimSpace(p.Name) == "" {
		return errors.New("invalid content preset")
	}
	var existing ContentPreset
	err := db.Where("name = ?", strings.TrimSpace(p.Name)).Where("role_id = ?", p.RoleID).First(&existing).Error
	if err == nil {
		p.ID = existing.ID
		existing.Name = p.Name
		existing.Language = p.Language
		existing.RoleID = p.RoleID
		existing.DataJSON = p.DataJSON
		existing.IsActive = p.IsActive
		return db.Save(&existing).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return db.Create(p).Error
	}
	return err
}

func (r *Repo) AdminListContentPresets(page, size int, q, role, language string) ([]ContentPreset, int64, error) {
	page = clampPage(page)
	size = clampPageSize(size)
	var list []ContentPreset
	db := r.db.Model(&ContentPreset{})
	if q = strings.TrimSpace(q); q != "" {
		qq := "%" + q + "%"
		db = db.Where("name LIKE ?", qq)
		if id, err := strconv.ParseUint(q, 10, 64); err == nil {
			db = db.Or("id = ?", uint(id))
		}
	}
	if role = strings.TrimSpace(role); role != "" {
		if roleID, err := strconv.ParseUint(role, 10, 64); err == nil {
			db = db.Where("role_id = ?", uint(roleID))
		}
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
	if p == nil || strings.TrimSpace(p.Name) == "" || p.RoleID == 0 {
		return errors.New("invalid")
	}
	if err := ValidatePresetDataJSON(p.DataJSON); err != nil {
		return err
	}
	return r.db.Create(p).Error
}

func (r *Repo) AdminPatchContentPreset(id uint, patch map[string]any) error {
	if v, ok := patch["data_json"]; ok {
		if s, ok := v.(string); ok {
			if err := ValidatePresetDataJSON(s); err != nil {
				return err
			}
		}
	}
	return r.db.Model(&ContentPreset{}).Where("id = ?", id).Updates(patch).Error
}

func (r *Repo) AdminDeleteContentPreset(id uint) error {
	return r.db.Where("id = ?", id).Delete(&ContentPreset{}).Error
}
