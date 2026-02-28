package stats

import (
	"time"

	"cvforge/internal/infra/database"
	"cvforge/internal/models"

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

func (r *Repo) CountUsers() (int64, error) {
	var n int64
	err := r.db.Model(&models.User{}).Count(&n).Error
	return n, err
}

func (r *Repo) CountResumes() (int64, error) {
	var n int64
	err := r.db.Model(&models.Resume{}).Count(&n).Error
	return n, err
}

func (r *Repo) CountTemplates() (int64, error) {
	var n int64
	err := r.db.Model(&models.Template{}).Count(&n).Error
	return n, err
}

func (r *Repo) CountUsersCreatedBetween(start, end time.Time) (int64, error) {
	var n int64
	err := r.db.Model(&models.User{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&n).Error
	return n, err
}

func (r *Repo) CountResumesCreatedBetween(start, end time.Time) (int64, error) {
	var n int64
	err := r.db.Model(&models.Resume{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&n).Error
	return n, err
}

func (r *Repo) CountTemplatesCreatedBetween(start, end time.Time) (int64, error) {
	var n int64
	err := r.db.Model(&models.Template{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&n).Error
	return n, err
}

