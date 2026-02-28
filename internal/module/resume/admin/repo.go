package admin

import (
	"errors"
	"time"

	"cvforge/internal/infra/database"
	resumemod "cvforge/internal/module/resume"

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

func (r *Repo) ListResumes(page, size int, userID *uint, titleLike string, templateID string) ([]resumemod.Resume, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	q := r.db.Model(&resumemod.Resume{}).Preload("Personal").Preload("Theme")
	if userID != nil && *userID != 0 {
		q = q.Where("user_id = ?", *userID)
	}
	if titleLike != "" {
		q = q.Where("title LIKE ?", "%"+titleLike+"%")
	}
	if templateID != "" {
		q = q.Where("template_id = ?", templateID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []resumemod.Resume
	err := q.Order("updated_at desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *Repo) FindUsersByIDs(ids []uint) ([]resumemod.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []resumemod.User
	err := r.db.Where("id IN ?", ids).Find(&list).Error
	return list, err
}

func (r *Repo) FindResumeByID(id uint) (resumemod.Resume, error) {
	var res resumemod.Resume
	err := r.db.Where("id = ?", id).First(&res).Error
	return res, err
}

func (r *Repo) FindResumeFullByID(id uint) (resumemod.Resume, error) {
	var res resumemod.Resume
	err := r.db.Where("id = ?", id).Preload("Personal").Preload("Theme").Preload("Sections.Items").First(&res).Error
	return res, err
}

func (r *Repo) DeleteResumeByID(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var secIDs []uint
		if err := tx.Model(&resumemod.ResumeSection{}).Where("resume_id = ?", id).Pluck("id", &secIDs).Error; err != nil {
			return err
		}
		if len(secIDs) > 0 {
			if err := tx.Where("section_id IN ?", secIDs).Delete(&resumemod.ResumeItem{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("resume_id = ?", id).Delete(&resumemod.ResumeSection{}).Error; err != nil {
			return err
		}
		if err := tx.Where("resume_id = ?", id).Delete(&resumemod.ResumePersonal{}).Error; err != nil {
			return err
		}
		if err := tx.Where("resume_id = ?", id).Delete(&resumemod.ResumeTheme{}).Error; err != nil {
			return err
		}
		if err := tx.Where("resume_id = ?", id).Delete(&resumemod.ShareLink{}).Error; err != nil {
			return err
		}
		return tx.Delete(&resumemod.Resume{Model: gorm.Model{ID: id}}).Error
	})
}

func (r *Repo) FindShareLinkByResumeID(resumeID uint) (resumemod.ShareLink, error) {
	var sl resumemod.ShareLink
	err := r.db.Where("resume_id = ?", resumeID).First(&sl).Error
	return sl, err
}

func (r *Repo) CreateShareLink(sl *resumemod.ShareLink) error {
	return r.db.Create(sl).Error
}

func (r *Repo) SaveShareLink(sl *resumemod.ShareLink) error {
	return r.db.Save(sl).Error
}

func (r *Repo) CreateAuditLog(al *resumemod.AuditLog) error {
	if al == nil {
		return errors.New("invalid")
	}
	return r.db.Create(al).Error
}

func (r *Repo) Now() time.Time {
	return time.Now()
}

