package models

import (
	"time"

	"gorm.io/gorm"
)

type ShareLink struct {
	gorm.Model
	ResumeID     uint   // 关联简历 ID
	UserID       uint   // 关联用户 ID
	Slug         string `gorm:"uniqueIndex;size:64"` // 分享链接唯一标识
	IsPublic     bool   // 是否公开可访问
	Password     string
	ExpiresAt    *time.Time
	Views        uint64
	LastAccessAt *time.Time
}

func (ShareLink) TableName() string {
	return "share_link"
}
