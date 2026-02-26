package taxonomy

import "cvforge/internal/models"

type JobCategory = models.JobCategory
type JobCategoryI18n = models.JobCategoryI18n
type JobRole = models.JobRole
type JobRoleI18n = models.JobRoleI18n

type JobCategoryView struct {
	JobCategory
	Name string `gorm:"column:name"`
}

type JobRoleView struct {
	JobRole
	Name string `gorm:"column:name"`
}
