package preset

import (
	"errors"
	"strconv"
	"strings"

	"openresume/internal/infra/database"
	"openresume/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo struct {
	db *gorm.DB
}

func normalizeLanguage(raw string) string {
	v := strings.TrimSpace(strings.ToLower(raw))
	if v == "en" {
		return "en"
	}
	if v == "zh" {
		return "zh"
	}
	return "zh"
}

func DefaultRepo() *Repo {
	return &Repo{db: database.DB}
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) getByIDActiveWithLanguage(id uint, language string) (ContentPreset, error) {
	var p ContentPreset
	err := r.db.Table("content_preset cp").
		Select("cp.id, cp.name, cp.role_id, cp.is_active, cpi.language, cpi.data_json").
		Joins("JOIN content_preset_i18n cpi ON cpi.content_preset_id = cp.id AND cpi.deleted_at IS NULL").
		Where("cp.id = ?", id).
		Where("cp.is_active = ?", true).
		Where("cp.deleted_at IS NULL").
		Where("cpi.language = ?", language).
		First(&p).Error
	return p, err
}

func (r *Repo) GetByIDActive(id uint, language string) (ContentPreset, error) {
	language = normalizeLanguage(language)

	p, err := r.getByIDActiveWithLanguage(id, language)
	if err == nil {
		return p, nil
	}
	p, err = r.getByIDActiveWithLanguage(id, "zh")
	if err == nil {
		return p, nil
	}

	var any ContentPreset
	err = r.db.Table("content_preset cp").
		Select("cp.id, cp.name, cp.role_id, cp.is_active, cpi.language, cpi.data_json").
		Joins("JOIN content_preset_i18n cpi ON cpi.content_preset_id = cp.id AND cpi.deleted_at IS NULL").
		Where("cp.id = ?", id).
		Where("cp.is_active = ?", true).
		Where("cp.deleted_at IS NULL").
		Order("cpi.id asc").
		Limit(1).
		First(&any).Error
	return any, err
}

func (r *Repo) UpsertContentPreset(db *gorm.DB, p *ContentPreset) error {
	if p == nil || strings.TrimSpace(p.Name) == "" || p.RoleID == 0 {
		return errors.New("invalid content preset")
	}
	if err := ValidatePresetDataJSON(p.DataJSON); err != nil {
		return err
	}

	var base models.ContentPreset
	err := db.Model(&models.ContentPreset{}).
		Where("name = ?", strings.TrimSpace(p.Name)).
		Where("role_id = ?", p.RoleID).
		First(&base).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		base = models.ContentPreset{
			Name:     strings.TrimSpace(p.Name),
			RoleID:   p.RoleID,
			IsActive: p.IsActive,
		}
		if err := db.Create(&base).Error; err != nil {
			return err
		}
	} else {
		base.Name = strings.TrimSpace(p.Name)
		base.RoleID = p.RoleID
		base.IsActive = p.IsActive
		if err := db.Save(&base).Error; err != nil {
			return err
		}
	}

	lang := normalizeLanguage(p.Language)
	data := strings.TrimSpace(p.DataJSON)
	i18n := models.ContentPresetI18n{
		ContentPresetID: base.ID,
		Language:        lang,
		DataJSON:        data,
	}
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "content_preset_id"}, {Name: "language"}},
		DoUpdates: clause.AssignmentColumns([]string{"data_json"}),
	}).Create(&i18n).Error; err != nil {
		return err
	}

	p.ID = base.ID
	p.Language = lang
	p.DataJSON = data
	p.IsActive = base.IsActive
	return nil
}

func (r *Repo) AdminListContentPresets(page, size int, q, role, language string) ([]ContentPreset, int64, error) {
	page = clampPage(page)
	size = clampPageSize(size)
	var list []ContentPreset
	language = normalizeLanguage(language)
	db := r.db.Table("content_preset cp").
		Select("cp.id, cp.name, cp.role_id, cp.is_active, cpi.language, cpi.data_json").
		Joins("JOIN content_preset_i18n cpi ON cpi.content_preset_id = cp.id AND cpi.deleted_at IS NULL AND cpi.language = ?", language).
		Where("cp.deleted_at IS NULL")
	if q = strings.TrimSpace(q); q != "" {
		qq := "%" + q + "%"
		if id, err := strconv.ParseUint(q, 10, 64); err == nil {
			db = db.Where("(cp.name LIKE ? OR cp.id = ?)", qq, uint(id))
		} else {
			db = db.Where("cp.name LIKE ?", qq)
		}
	}
	if role = strings.TrimSpace(role); role != "" {
		if roleID, err := strconv.ParseUint(role, 10, 64); err == nil {
			db = db.Where("cp.role_id = ?", uint(roleID))
		}
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := db.Order("cp.id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *Repo) AdminCreateContentPreset(p *ContentPreset) error {
	if p == nil || strings.TrimSpace(p.Name) == "" || p.RoleID == 0 {
		return errors.New("invalid")
	}
	if err := ValidatePresetDataJSON(p.DataJSON); err != nil {
		return err
	}
	name := strings.TrimSpace(p.Name)
	lang := normalizeLanguage(p.Language)
	data := strings.TrimSpace(p.DataJSON)

	return r.db.Transaction(func(tx *gorm.DB) error {
		var base models.ContentPreset
		err := tx.Model(&models.ContentPreset{}).
			Where("name = ?", name).
			Where("role_id = ?", p.RoleID).
			First(&base).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			base = models.ContentPreset{
				Name:     name,
				RoleID:   p.RoleID,
				IsActive: p.IsActive,
			}
			if err := tx.Create(&base).Error; err != nil {
				return err
			}
		} else {
			base.IsActive = p.IsActive
			if err := tx.Save(&base).Error; err != nil {
				return err
			}
		}
		i18n := models.ContentPresetI18n{ContentPresetID: base.ID, Language: lang, DataJSON: data}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "content_preset_id"}, {Name: "language"}},
			DoUpdates: clause.AssignmentColumns([]string{"data_json"}),
		}).Create(&i18n).Error; err != nil {
			return err
		}
		p.ID = base.ID
		p.Language = lang
		p.DataJSON = data
		p.IsActive = base.IsActive
		return nil
	})
}

func (r *Repo) AdminPatchContentPreset(id uint, patch map[string]any) error {
	var basePatch = map[string]any{}
	var targetLanguage = "zh"
	var dataJSON *string

	if v, ok := patch["name"]; ok {
		basePatch["name"] = strings.TrimSpace(v.(string))
	}
	if v, ok := patch["role_id"]; ok {
		basePatch["role_id"] = v
	}
	if v, ok := patch["is_active"]; ok {
		basePatch["is_active"] = v
	}
	if v, ok := patch["language"]; ok {
		if s, ok := v.(string); ok {
			targetLanguage = normalizeLanguage(s)
		}
	}
	if v, ok := patch["data_json"]; ok {
		if s, ok := v.(string); ok {
			if err := ValidatePresetDataJSON(s); err != nil {
				return err
			}
			ss := strings.TrimSpace(s)
			dataJSON = &ss
		}
	}
	if len(basePatch) == 0 && dataJSON == nil {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if len(basePatch) != 0 {
			if err := tx.Model(&models.ContentPreset{}).Where("id = ?", id).Updates(basePatch).Error; err != nil {
				return err
			}
		}
		if dataJSON != nil {
			i18n := models.ContentPresetI18n{
				ContentPresetID: id,
				Language:        targetLanguage,
				DataJSON:        *dataJSON,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "content_preset_id"}, {Name: "language"}},
				DoUpdates: clause.AssignmentColumns([]string{"data_json"}),
			}).Create(&i18n).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repo) AdminDeleteContentPreset(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&models.ContentPreset{}).Error; err != nil {
			return err
		}
		if err := tx.Where("content_preset_id = ?", id).Delete(&models.ContentPresetI18n{}).Error; err != nil {
			return err
		}
		return nil
	})
}
