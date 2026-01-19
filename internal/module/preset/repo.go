package preset

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

func (r *Repo) GetByExternalActive(id string) (ContentPreset, error) {
	var p ContentPreset
	err := r.db.Where("external_id = ?", id).Where("is_active = ?", true).First(&p).Error
	return p, err
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
	if err := ValidatePresetDataJSON(p.DataJSON); err != nil {
		return err
	}
	return r.db.Create(p).Error
}

func (r *Repo) AdminPatchContentPreset(externalID string, patch map[string]any) error {
	if v, ok := patch["data_json"]; ok {
		if s, ok := v.(string); ok {
			if err := ValidatePresetDataJSON(s); err != nil {
				return err
			}
		}
	}
	return r.db.Model(&ContentPreset{}).Where("external_id = ?", externalID).Updates(patch).Error
}

func (r *Repo) AdminDeleteContentPreset(externalID string) error {
	return r.db.Where("external_id = ?", externalID).Delete(&ContentPreset{}).Error
}
