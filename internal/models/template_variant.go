package models

import (
	"time"

	"gorm.io/gorm"
)

type TemplateVariant struct {
	ID                      uint `gorm:"primarykey"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               gorm.DeletedAt `gorm:"uniqueIndex:uk_external_id_deleted_at,priority:10"`
	ExternalID              string         `gorm:"uniqueIndex:uk_external_id_deleted_at,priority:1;size:64"`
	Name                    string         `gorm:"size:128"`
	LayoutTemplateExternalID string        `gorm:"size:64;index"`
	PresetExternalID        string         `gorm:"size:64;index"`
	RoleExternalID          string         `gorm:"size:64;index"`
	Tags                    string         `gorm:"size:256"` // comma-separated
	UsageCount              int            `gorm:"default:0"`
	IsPremium               bool           `gorm:"default:false"`
	IsActive                bool           `gorm:"default:true"`
}

func (TemplateVariant) TableName() string {
	return "template_variant"
}

