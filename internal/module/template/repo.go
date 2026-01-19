package template

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

func (r *Repo) Count() (int64, error) {
	var count int64
	err := r.db.Model(&Template{}).Count(&count).Error
	return count, err
}

func (r *Repo) ListAll() ([]Template, error) {
	var list []Template
	err := r.db.Order("name asc").Order("id asc").Find(&list).Error
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
