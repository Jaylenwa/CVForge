package models

import "gorm.io/gorm"

type JobRoleI18n struct {
	gorm.Model
	JobRoleID uint           `gorm:"index;index:uniq_job_role_language,unique"`
	Language  string         `gorm:"size:16;index:uniq_job_role_language,unique"`
	Name      string         `gorm:"size:128"`
}

func (JobRoleI18n) TableName() string {
	return "job_role_i18n"
}
