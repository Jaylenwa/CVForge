package models

import "gorm.io/gorm"

type ShareLink struct {
	gorm.Model
	ResumeID uint
	Slug     string `gorm:"uniqueIndex;size:64"`
	IsPublic bool
}

func (ShareLink) TableName() string {
	return "share_link"
}
