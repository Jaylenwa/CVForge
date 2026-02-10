package models

import "gorm.io/gorm"

type JobCategoryI18n struct {
	gorm.Model
	JobCategoryID uint   `gorm:"index;index:uniq_job_category_language,unique"`
	Language      string `gorm:"size:16;index:uniq_job_category_language,unique"`
	Name          string `gorm:"size:128"`
}

func (JobCategoryI18n) TableName() string {
	return "job_category_i18n"
}
