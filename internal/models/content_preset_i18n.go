package models

import "gorm.io/gorm"

type ContentPresetI18n struct {
	gorm.Model
	ContentPresetID uint           `gorm:"index;index:uniq_content_preset_language,unique"`
	Language        string         `gorm:"size:16;index;index:uniq_content_preset_language,unique"`
	DataJSON        string         `gorm:"type:longtext"`
}

func (ContentPresetI18n) TableName() string {
	return "content_preset_i18n"
}
