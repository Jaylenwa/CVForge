package share

import (
	"time"

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

func (r *Repo) FindResumeByID(resumeID uint) (Resume, error) {
	var res Resume
	err := r.db.Where("id = ?", resumeID).First(&res).Error
	return res, err
}

func (r *Repo) FindResumeWithPublicPreloadsByID(resumeID uint) (Resume, error) {
	var res Resume
	err := r.db.Where("id = ?", resumeID).
		Preload("Personal").
		Preload("Theme").
		Preload("Sections.Items").
		First(&res).Error
	return res, err
}

func (r *Repo) UpdateViewsAndLastAccess(id uint, views uint64, lastAccessAt time.Time) error {
	return r.db.Model(&ShareLink{}).Where("id = ?", id).Updates(map[string]any{
		"views":          views,
		"last_access_at": lastAccessAt,
	}).Error
}

func (r *Repo) IncrementViewsAndLastAccess(id uint, lastAccessAt time.Time) error {
	return r.db.Model(&ShareLink{}).Where("id = ?", id).Updates(map[string]any{
		"views":          gorm.Expr("views + 1"),
		"last_access_at": lastAccessAt,
	}).Error
}
