package database

import (
	"errors"
	"strings"

	"openresume/internal/models"

	"gorm.io/gorm"
)

const migrationExternalIDToIDDoneKey = "migration_externalid_to_id_done"

func RunExternalIDToIDBackfill(db *gorm.DB) error {
	if db == nil {
		return errors.New("db nil")
	}
	var cfg models.Config
	err := db.Where("config_key = ?", migrationExternalIDToIDDoneKey).First(&cfg).Error
	if err == nil && strings.TrimSpace(cfg.ConfigValue) == "done" {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var legacyProbe any
		if err := tx.Raw(`SELECT external_id FROM job_role LIMIT 1`).Scan(&legacyProbe).Error; err != nil {
			var existing models.Config
			if err := tx.Where("config_key = ?", migrationExternalIDToIDDoneKey).First(&existing).Error; err == nil {
				return tx.Model(&models.Config{}).Where("id = ?", existing.ID).Update("config_value", "done").Error
			}
			return tx.Create(&models.Config{ConfigKey: migrationExternalIDToIDDoneKey, ConfigValue: "done", Description: "backfill external_id to numeric ids", Type: "string"}).Error
		}

		type CatRow struct {
			ID               uint
			ParentExternalID string
		}
		var cats []CatRow
		if err := tx.Raw(`SELECT id, parent_external_id FROM job_category WHERE (parent_id IS NULL OR parent_id = 0) AND parent_external_id IS NOT NULL AND parent_external_id <> '' AND deleted_at IS NULL`).Scan(&cats).Error; err != nil {
			return err
		}
		for _, c := range cats {
			var parentID uint
			if err := tx.Raw(`SELECT id FROM job_category WHERE external_id = ? AND deleted_at IS NULL LIMIT 1`, c.ParentExternalID).Scan(&parentID).Error; err != nil {
				return err
			}
			if parentID == 0 {
				continue
			}
			if err := tx.Exec(`UPDATE job_category SET parent_id = ? WHERE id = ?`, parentID, c.ID).Error; err != nil {
				return err
			}
		}

		type RoleRow struct {
			ID                 uint
			CategoryExternalID string
		}
		var roles []RoleRow
		if err := tx.Raw(`SELECT id, category_external_id FROM job_role WHERE (category_id IS NULL OR category_id = 0) AND category_external_id IS NOT NULL AND category_external_id <> '' AND deleted_at IS NULL`).Scan(&roles).Error; err != nil {
			return err
		}
		for _, r := range roles {
			var categoryID uint
			if err := tx.Raw(`SELECT id FROM job_category WHERE external_id = ? AND deleted_at IS NULL LIMIT 1`, r.CategoryExternalID).Scan(&categoryID).Error; err != nil {
				return err
			}
			if categoryID == 0 {
				continue
			}
			if err := tx.Exec(`UPDATE job_role SET category_id = ? WHERE id = ?`, categoryID, r.ID).Error; err != nil {
				return err
			}
		}

		type PresetRow struct {
			ID             uint
			RoleExternalID string
		}
		var presets []PresetRow
		if err := tx.Raw(`SELECT id, role_external_id FROM content_preset WHERE (role_id IS NULL OR role_id = 0) AND role_external_id IS NOT NULL AND role_external_id <> '' AND deleted_at IS NULL`).Scan(&presets).Error; err != nil {
			return err
		}
		for _, p := range presets {
			var roleID uint
			if err := tx.Raw(`SELECT id FROM job_role WHERE external_id = ? AND deleted_at IS NULL LIMIT 1`, p.RoleExternalID).Scan(&roleID).Error; err != nil {
				return err
			}
			if roleID == 0 {
				continue
			}
			if err := tx.Exec(`UPDATE content_preset SET role_id = ? WHERE id = ?`, roleID, p.ID).Error; err != nil {
				return err
			}
		}

		type VariantRow struct {
			ID               uint
			PresetExternalID string
			RoleExternalID   string
		}
		var variants []VariantRow
		if err := tx.Raw(`SELECT id, preset_external_id, role_external_id FROM template_variant WHERE (preset_id IS NULL OR preset_id = 0 OR role_id IS NULL OR role_id = 0) AND deleted_at IS NULL`).Scan(&variants).Error; err != nil {
			return err
		}
		for _, v := range variants {
			var presetID uint
			if v.PresetExternalID != "" {
				_ = tx.Raw(`SELECT id FROM content_preset WHERE external_id = ? AND deleted_at IS NULL LIMIT 1`, v.PresetExternalID).Scan(&presetID).Error
			}
			var roleID uint
			if v.RoleExternalID != "" {
				_ = tx.Raw(`SELECT id FROM job_role WHERE external_id = ? AND deleted_at IS NULL LIMIT 1`, v.RoleExternalID).Scan(&roleID).Error
			}
			if presetID != 0 {
				if err := tx.Exec(`UPDATE template_variant SET preset_id = ? WHERE id = ?`, presetID, v.ID).Error; err != nil {
					return err
				}
			}
			if roleID != 0 {
				if err := tx.Exec(`UPDATE template_variant SET role_id = ? WHERE id = ?`, roleID, v.ID).Error; err != nil {
					return err
				}
			}
		}

		type ResumeRow struct {
			ID              uint
			VariantExternal string
			PresetExternal  string
			RoleExternal    string
		}
		var resumes []ResumeRow
		if err := tx.Raw(`SELECT id, variant_id AS variant_external, preset_id AS preset_external, role_id AS role_external FROM resume WHERE deleted_at IS NULL`).Scan(&resumes).Error; err != nil {
			return err
		}
		for _, r := range resumes {
			if r.VariantExternal != "" {
				var vid uint
				_ = tx.Raw(`SELECT id FROM template_variant WHERE external_id = ? AND deleted_at IS NULL LIMIT 1`, r.VariantExternal).Scan(&vid).Error
				if vid != 0 {
					if err := tx.Exec(`UPDATE resume SET variant_db_id = ? WHERE id = ?`, vid, r.ID).Error; err != nil {
						return err
					}
				}
			}
			if r.PresetExternal != "" {
				var pid uint
				_ = tx.Raw(`SELECT id FROM content_preset WHERE external_id = ? AND deleted_at IS NULL LIMIT 1`, r.PresetExternal).Scan(&pid).Error
				if pid != 0 {
					if err := tx.Exec(`UPDATE resume SET preset_db_id = ? WHERE id = ?`, pid, r.ID).Error; err != nil {
						return err
					}
				}
			}
			if r.RoleExternal != "" {
				var rid uint
				_ = tx.Raw(`SELECT id FROM job_role WHERE external_id = ? AND deleted_at IS NULL LIMIT 1`, r.RoleExternal).Scan(&rid).Error
				if rid != 0 {
					if err := tx.Exec(`UPDATE resume SET role_db_id = ? WHERE id = ?`, rid, r.ID).Error; err != nil {
						return err
					}
				}
			}
		}

		var existing models.Config
		if err := tx.Where("config_key = ?", migrationExternalIDToIDDoneKey).First(&existing).Error; err == nil {
			return tx.Model(&models.Config{}).Where("id = ?", existing.ID).Update("config_value", "done").Error
		}
		return tx.Create(&models.Config{ConfigKey: migrationExternalIDToIDDoneKey, ConfigValue: "done", Description: "backfill external_id to numeric ids", Type: "string"}).Error
	})
}
