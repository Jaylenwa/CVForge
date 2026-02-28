package admin

import (
	"cvforge/internal/infra/database"
	sharemod "cvforge/internal/module/share"

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

func (r *Repo) FindShareLinkBySlug(slug string) (sharemod.ShareLink, error) {
	var sl sharemod.ShareLink
	err := r.db.Where("slug = ?", slug).First(&sl).Error
	return sl, err
}

func (r *Repo) SaveShareLink(sl *sharemod.ShareLink) error {
	return r.db.Save(sl).Error
}

func (r *Repo) DeleteBySlug(slug string) error {
	return r.db.Where("slug = ?", slug).Delete(&sharemod.ShareLink{}).Error
}

func (r *Repo) AdminListShareLinks(page, size int, slug string, isPublic *bool) ([]sharemod.ShareLink, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	var list []sharemod.ShareLink
	q := r.db.Model(&sharemod.ShareLink{})
	if slug != "" {
		q = q.Where("slug LIKE ?", "%"+slug+"%")
	}
	if isPublic != nil {
		q = q.Where("is_public = ?", *isPublic)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *Repo) FindResumesByIDs(ids []uint) ([]sharemod.Resume, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []sharemod.Resume
	err := r.db.Where("id IN ?", ids).Find(&list).Error
	return list, err
}

func (r *Repo) FindUsersByIDs(ids []uint) ([]sharemod.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []sharemod.User
	err := r.db.Where("id IN ?", ids).Find(&list).Error
	return list, err
}

