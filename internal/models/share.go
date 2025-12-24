package models

import "gorm.io/gorm"

type ShareLink struct {
	gorm.Model
	ResumeID uint
	UserID   uint
	Slug     string `gorm:"uniqueIndex;size:64"`
	IsPublic bool
}

func (ShareLink) TableName() string {
	return "share_link"
}
