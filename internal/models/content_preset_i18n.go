package models

import (
	"time"

	"gorm.io/gorm"
)

type ContentPresetI18n struct {
	ID              uint `gorm:"primarykey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	ContentPresetID uint           `gorm:"index;index:uniq_content_preset_language,unique"`
	Language        string         `gorm:"size:16;index;index:uniq_content_preset_language,unique"`
	DataJSON        string         `gorm:"type:longtext"`
}

func (ContentPresetI18n) TableName() string {
	return "content_preset_i18n"
}
