package user

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

func (r *Repo) FindByID(id any) (User, error) {
	var u User
	err := r.db.First(&u, id).Error
	return u, err
}

func (r *Repo) Save(u *User) error {
	return r.db.Save(u).Error
}
