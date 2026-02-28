package pdf

import (
	"cvforge/internal/infra/database"

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

func (r *Repo) FindResumeWithSections(id uint) (Resume, error) {
	var res Resume
	err := r.db.Where("id = ?", id).Preload("Sections.Items").First(&res).Error
	return res, err
}

func (r *Repo) FindResumeTitle(id uint) (string, error) {
	var res Resume
	if err := r.db.Select("id, title").Where("id = ?", id).First(&res).Error; err != nil {
		return "", err
	}
	return res.Title, nil
}

