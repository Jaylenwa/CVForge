package models

import (
	"time"

	"gorm.io/gorm"
)

type ContentPreset struct {
	ID             uint           `gorm:"primarykey"` // 主键
	CreatedAt      time.Time      // 创建时间
	UpdatedAt      time.Time      // 更新时间
	DeletedAt      gorm.DeletedAt `gorm:"index"` // 软删除时间
	Name           string         `gorm:"size:128"`                                                                // 预设名称（用于展示）
	Language       string         `gorm:"size:8"`                                                                  // 语言（zh/en）
	RoleID         uint           `gorm:"index"`                                                                   // 关联岗位方向 ID
	Tags           string         `gorm:"size:256"`                                                                // 标签（逗号分隔）
	DataJSON       string         `gorm:"type:longtext"`                                                           // 预设内容（JSON 字符串，包含个人信息/主题/sections）
	IsActive       bool           `gorm:"default:true"`                                                            // 是否启用
}

func (ContentPreset) TableName() string {
	return "content_preset"
}
