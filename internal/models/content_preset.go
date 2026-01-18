package models

import (
	"time"

	"gorm.io/gorm"
)

type ContentPreset struct {
	ID             uint           `gorm:"primarykey"` // 主键
	CreatedAt      time.Time      // 创建时间
	UpdatedAt      time.Time      // 更新时间
	DeletedAt      gorm.DeletedAt `gorm:"uniqueIndex:uk_content_preset_external_id_deleted_at,priority:10"`        // 软删除时间（与 ExternalID 共同唯一）
	ExternalID     string         `gorm:"uniqueIndex:uk_content_preset_external_id_deleted_at,priority:1;size:64"` // 业务外部 ID（内容预设唯一标识）
	Name           string         `gorm:"size:128"`                                                                // 预设名称（用于展示）
	Language       string         `gorm:"size:8"`                                                                  // 语言（zh/en）
	RoleExternalID string         `gorm:"size:64;index"`                                                           // 关联岗位方向 ExternalID
	Tags           string         `gorm:"size:256"`                                                                // 标签（逗号分隔）
	DataJSON       string         `gorm:"type:longtext"`                                                           // 预设内容（JSON 字符串，包含个人信息/主题/sections）
	IsActive       bool           `gorm:"default:true"`                                                            // 是否启用
}

func (ContentPreset) TableName() string {
	return "content_preset"
}
