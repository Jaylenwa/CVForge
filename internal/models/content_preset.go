package models

import (
	"time"

	"gorm.io/gorm"
)

type ContentPreset struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"uniqueIndex:uk_external_id_deleted_at,priority:10"`
	ExternalID   string         `gorm:"uniqueIndex:uk_external_id_deleted_at,priority:1;size:64"`
	Name         string         `gorm:"size:128"`
	Language     string         `gorm:"size:8"`
	RoleExternalID string       `gorm:"size:64;index"`
	Tags         string         `gorm:"size:256"` // comma-separated
	DataJSON     string         `gorm:"type:longtext"`
	IsActive     bool           `gorm:"default:true"`
}

func (ContentPreset) TableName() string {
	return "content_preset"
}

