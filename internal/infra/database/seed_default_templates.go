package database

import (
	"errors"
	"strings"

	"openresume/internal/models"

	"gorm.io/gorm"
)

const seedDefaultTemplatesDoneKey = "seed_default_templates_v2"

func RunDefaultTemplatesSeed(db *gorm.DB) error {
	if db == nil {
		return errors.New("db nil")
	}
	var cfg models.Config
	err := db.Where("config_key = ?", seedDefaultTemplatesDoneKey).First(&cfg).Error
	if err == nil && strings.TrimSpace(cfg.ConfigValue) == "done" {
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	type item struct {
		ExternalID string
		Name       string
		Tags       string
	}
	defaults := []item{
		{ExternalID: "TemplateClassic", Name: "经典专业版", Tags: "专业,简洁,ATS 友好"},
		{ExternalID: "TemplateMintTimeline", Name: "青色时间轴", Tags: "美观,中文,ATS 友好"},
		{ExternalID: "TemplateSlateSidebar", Name: "石板简洁版", Tags: "单栏,简洁,ATS 友好"},
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, it := range defaults {
			var existing models.Template
			err := tx.Unscoped().Where("external_id = ?", it.ExternalID).First(&existing).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if err == nil {
				if existing.DeletedAt.Valid {
					continue
				}
				if it.ExternalID == "TemplateSlateSidebar" && existing.Name == "侧边栏双栏" && existing.Tags == "双栏,ATS 友好,简洁" {
					if err := tx.Model(&models.Template{}).Where("id = ?", existing.ID).Updates(map[string]any{"name": it.Name, "tags": it.Tags}).Error; err != nil {
						return err
					}
				}
				continue
			}
			if err := tx.Create(&models.Template{ExternalID: it.ExternalID, Name: it.Name, Tags: it.Tags}).Error; err != nil {
				return err
			}
		}

		var cfg models.Config
		if err := tx.Where("config_key = ?", seedDefaultTemplatesDoneKey).First(&cfg).Error; err == nil {
			return tx.Model(&models.Config{}).Where("id = ?", cfg.ID).Update("config_value", "done").Error
		}
		return tx.Create(&models.Config{ConfigKey: seedDefaultTemplatesDoneKey, ConfigValue: "done", Description: "seed default templates into template table", Type: "string"}).Error
	})
}
