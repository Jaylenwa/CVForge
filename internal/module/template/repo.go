package template

import (
	"errors"
	"strings"

	"openresume/internal/infra/database"
	"openresume/internal/models"

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

func (r *Repo) Count() (int64, error) {
	var count int64
	err := r.db.Model(&Template{}).Count(&count).Error
	return count, err
}

func (r *Repo) ListAll() ([]Template, error) {
	var list []Template
	err := r.db.Order("id asc").Find(&list).Error
	return list, err
}

func (r *Repo) GetByExternal(id string) (Template, error) {
	var t Template
	err := r.db.Where("external_id = ?", id).First(&t).Error
	return t, err
}

func (r *Repo) Create(t *Template) error {
	return r.db.Create(t).Error
}

func (r *Repo) Save(t *Template) error {
	return r.db.Save(t).Error
}

func (r *Repo) DeleteByExternal(id string) error {
	return r.db.Where("external_id = ?", id).Delete(&Template{}).Error
}

func (r *Repo) ListI18n(templateIDs []uint) ([]models.TemplateI18n, error) {
	if len(templateIDs) == 0 {
		return nil, nil
	}
	var list []models.TemplateI18n
	err := r.db.Where("template_id IN ?", templateIDs).Find(&list).Error
	return list, err
}

func upsertTemplateI18n(tx *gorm.DB, templateID uint, language, name string) error {
	language = strings.ToLower(strings.TrimSpace(language))
	if language != "zh" && language != "en" {
		return nil
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	m := models.TemplateI18n{
		TemplateID: templateID,
		Language:   language,
		Name:       name,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "template_id"}, {Name: "language"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&m).Error
}

func (r *Repo) CreateWithNames(t *Template, names map[string]string) error {
	if t == nil || strings.TrimSpace(t.ExternalID) == "" {
		return errors.New("invalid")
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(t).Error; err != nil {
			return err
		}
		for lang, name := range names {
			if err := upsertTemplateI18n(tx, t.ID, lang, name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repo) PatchWithNames(templateID uint, patch map[string]any, names map[string]string) error {
	if templateID == 0 {
		return errors.New("invalid")
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if len(patch) > 0 {
			if err := tx.Model(&Template{}).Where("id = ?", templateID).Updates(patch).Error; err != nil {
				return err
			}
		}
		for lang, name := range names {
			if err := upsertTemplateI18n(tx, templateID, lang, name); err != nil {
				return err
			}
		}
		return nil
	})
}
