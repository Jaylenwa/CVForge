package share

import (
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

func (r *Repo) FindByResumeID(resumeID uint) (ShareLink, error) {
	var sl ShareLink
	err := r.db.Where("resume_id = ?", resumeID).First(&sl).Error
	return sl, err
}

func (r *Repo) FindPublicBySlug(slug string) (ShareLink, error) {
	var sl ShareLink
	err := r.db.Where("slug = ? AND is_public = ?", slug, true).First(&sl).Error
	return sl, err
}

func (r *Repo) FindBySlug(slug string) (ShareLink, error) {
	var sl ShareLink
	err := r.db.Where("slug = ?", slug).First(&sl).Error
	return sl, err
}

func (r *Repo) Create(sl *ShareLink) error {
	return r.db.Create(sl).Error
}

func (r *Repo) Save(sl *ShareLink) error {
	return r.db.Save(sl).Error
}
