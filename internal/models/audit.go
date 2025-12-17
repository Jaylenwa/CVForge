package models

import "gorm.io/gorm"

type AuditLog struct {
	gorm.Model
	ActorID    uint
	Action     string `gorm:"size:64"`
	TargetType string `gorm:"size:32"`
	TargetID   string `gorm:"size:128"`
	Metadata   string `gorm:"type:text"`
	IP         string `gorm:"size:64"`
	UA         string `gorm:"size:256"`
}

