package models

import (
	"time"
)

type RoleTemplateUsage struct {
	ID                uint      `gorm:"primarykey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	RoleID            uint      `gorm:"index;uniqueIndex:uk_role_template_usage_role_tpl,priority:1"`
	TemplateExternalID string    `gorm:"size:64;index;uniqueIndex:uk_role_template_usage_role_tpl,priority:2"`
	UsageCount        int       `gorm:"default:0"`
	LastUsedAt        *time.Time `gorm:"index"`
}

func (RoleTemplateUsage) TableName() string {
	return "role_template_usage"
}
