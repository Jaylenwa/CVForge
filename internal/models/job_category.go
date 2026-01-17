package models

import (
	"time"

	"gorm.io/gorm"
)

type JobCategory struct {
	ID               uint `gorm:"primarykey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"uniqueIndex:uk_external_id_deleted_at,priority:10"`
	ExternalID       string         `gorm:"uniqueIndex:uk_external_id_deleted_at,priority:1;size:64"`
	Name             string         `gorm:"size:128"`
	ParentExternalID string         `gorm:"size:64;index"`
	OrderNum         int            `gorm:"default:0"`
	IsActive         bool           `gorm:"default:true"`
}

func (JobCategory) TableName() string {
	return "job_category"
}

