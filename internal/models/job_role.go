package models

import (
	"time"

	"gorm.io/gorm"
)

type JobRole struct {
	ID                 uint `gorm:"primarykey"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"uniqueIndex:uk_external_id_deleted_at,priority:10"`
	ExternalID         string         `gorm:"uniqueIndex:uk_external_id_deleted_at,priority:1;size:64"`
	CategoryExternalID string         `gorm:"size:64;index"`
	Name               string         `gorm:"size:128"`
	Tags               string         `gorm:"size:256"` // comma-separated
	OrderNum           int            `gorm:"default:0"`
	IsActive           bool           `gorm:"default:true"`
}

func (JobRole) TableName() string {
	return "job_role"
}

