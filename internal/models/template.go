package models

import (
	"time"

	"gorm.io/gorm"
)

type Template struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"uniqueIndex:uk_external_id_deleted_at,priority:10"`
	ExternalID string         `gorm:"uniqueIndex:uk_external_id_deleted_at,priority:1;size:64"`
	Name       string         `gorm:"size:128"`
	Tags       string         `gorm:"size:256"` // comma-separated for简单实现
	UsageCount int            `gorm:"default:0"`
	IsPremium  bool           `gorm:"default:false"`
	Category   string         `gorm:"size:64"`
}

func (Template) TableName() string {
	return "template"
}
