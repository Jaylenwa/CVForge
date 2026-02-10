package models

import "gorm.io/gorm"

type TemplateI18n struct {
	gorm.Model
	TemplateID uint   `gorm:"index;index:uniq_template_language,unique"`
	Language   string `gorm:"size:16;index:uniq_template_language,unique"`
	Name       string `gorm:"size:128"`
}

func (TemplateI18n) TableName() string {
	return "template_i18n"
}
