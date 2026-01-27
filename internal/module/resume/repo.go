package resume

import (
	"time"

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

func (r *Repo) ListByUser(uid uint) ([]Resume, error) {
	var list []Resume
	err := r.db.Where("user_id = ?", uid).
		Preload("Personal").
		Preload("Theme").
		Preload("Sections.Items").
		Order("updated_at desc").
		Find(&list).Error
	return list, err
}

func (r *Repo) FindByID(id uint, preload bool) (Resume, error) {
	var res Resume
	q := r.db.Where("id = ?", id)
	if preload {
		q = q.Preload("Personal").Preload("Theme").Preload("Sections.Items")
	}
	err := q.First(&res).Error
	return res, err
}

func (r *Repo) Create(res *Resume) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(res).Error
}

func (r *Repo) DeleteByID(id uint) error {
	return r.db.Delete(&Resume{Model: gorm.Model{ID: id}}).Error
}

func (r *Repo) Replace(existing Resume, updated Resume) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var secIDs []uint
		if err := tx.Model(&ResumeSection{}).Where("resume_id = ?", existing.ID).Pluck("id", &secIDs).Error; err != nil {
			return err
		}
		if len(secIDs) > 0 {
			if err := tx.Unscoped().Where("section_id IN ?", secIDs).Delete(&ResumeItem{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Unscoped().Where("resume_id = ?", existing.ID).Delete(&ResumeSection{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("resume_id = ?", existing.ID).Delete(&ResumePersonal{}).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Where("resume_id = ?", existing.ID).Delete(&ResumeTheme{}).Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Omit("CreatedAt").Save(&updated).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *Repo) IncrementTemplateUsage(externalID string) error {
	if externalID == "" {
		return nil
	}
	return r.db.Model(&models.Template{}).
		Where("external_id = ?", externalID).
		UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}

func (r *Repo) IncrementRoleTemplateUsage(roleID uint, templateExternalID string) error {
	if roleID == 0 || templateExternalID == "" {
		return nil
	}
	now := time.Now()
	row := models.RoleTemplateUsage{
		RoleID:             roleID,
		TemplateExternalID: templateExternalID,
		UsageCount:         1,
		LastUsedAt:         &now,
	}
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "role_id"},
			{Name: "template_external_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"usage_count":  gorm.Expr("usage_count + ?", 1),
			"updated_at":   now,
			"last_used_at": now,
		}),
	}).Create(&row).Error
}
