package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;size:191"`
	PasswordHash string `gorm:"size:255"`
	Name         string `gorm:"size:128"`
	AvatarURL    string `gorm:"size:512"`
	Language     string `gorm:"size:8"`
	Role         string `gorm:"size:16;default:user"`
	IsActive     bool   `gorm:"default:true"`
	LastLoginAt  *time.Time
}

func (User) TableName() string {
	return "user"
}
