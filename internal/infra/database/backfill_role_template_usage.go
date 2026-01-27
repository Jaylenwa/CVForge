package database

import (
	"errors"
	"strings"
	"time"

	"openresume/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const migrationRoleTemplateUsageBackfillDoneKey = "migration_role_template_usage_backfill_done"

func RunRoleTemplateUsageBackfill(db *gorm.DB) error {
	if db == nil {
		return errors.New("db nil")
	}
	var cfg models.Config
	err := db.Where("config_key = ?", migrationRoleTemplateUsageBackfillDoneKey).First(&cfg).Error
	if err == nil && strings.TrimSpace(cfg.ConfigValue) == "done" {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		type AggRow struct {
			RoleID            uint      `gorm:"column:role_id"`
			TemplateExternalID string    `gorm:"column:template_external_id"`
			UsageCount        int       `gorm:"column:usage_count"`
			LastUsedAt        time.Time `gorm:"column:last_used_at"`
		}
		var rows []AggRow
		if err := tx.Raw(`
			SELECT
				role_db_id AS role_id,
				template_id AS template_external_id,
				COUNT(*) AS usage_count,
				MAX(updated_at) AS last_used_at
			FROM resume
			WHERE deleted_at IS NULL
			  AND role_db_id IS NOT NULL
			  AND role_db_id <> 0
			  AND template_id IS NOT NULL
			  AND template_id <> ''
			GROUP BY role_db_id, template_id
		`).Scan(&rows).Error; err != nil {
			return err
		}

		usages := make([]models.RoleTemplateUsage, 0, len(rows))
		now := time.Now()
		for _, r := range rows {
			last := r.LastUsedAt
			usages = append(usages, models.RoleTemplateUsage{
				RoleID:             r.RoleID,
				TemplateExternalID: r.TemplateExternalID,
				UsageCount:         r.UsageCount,
				LastUsedAt:         &last,
				UpdatedAt:          now,
			})
		}
		if len(usages) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "role_id"},
					{Name: "template_external_id"},
				},
				DoUpdates: clause.AssignmentColumns([]string{"usage_count", "updated_at", "last_used_at"}),
			}).CreateInBatches(&usages, 500).Error; err != nil {
				return err
			}
		}

		var maxTemplateUsage int64
		_ = tx.Model(&models.Template{}).Select("MAX(usage_count)").Scan(&maxTemplateUsage).Error
		if maxTemplateUsage == 0 {
			type TRow struct {
				TemplateExternalID string `gorm:"column:template_external_id"`
				UsageCount         int    `gorm:"column:usage_count"`
			}
			var trows []TRow
			_ = tx.Raw(`
				SELECT template_id AS template_external_id, COUNT(*) AS usage_count
				FROM resume
				WHERE deleted_at IS NULL
				  AND template_id IS NOT NULL
				  AND template_id <> ''
				GROUP BY template_id
			`).Scan(&trows).Error
			for _, tr := range trows {
				_ = tx.Model(&models.Template{}).
					Where("external_id = ?", tr.TemplateExternalID).
					UpdateColumn("usage_count", tr.UsageCount).Error
			}
		}

		var existing models.Config
		if err := tx.Where("config_key = ?", migrationRoleTemplateUsageBackfillDoneKey).First(&existing).Error; err == nil {
			return tx.Model(&models.Config{}).Where("id = ?", existing.ID).Update("config_value", "done").Error
		}
		return tx.Create(&models.Config{ConfigKey: migrationRoleTemplateUsageBackfillDoneKey, ConfigValue: "done", Description: "backfill role_template_usage and template usage_count from resume", Type: "string"}).Error
	})
}
