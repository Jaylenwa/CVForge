package models

import (
	"time"

	"gorm.io/gorm"
)

type ContentPreset struct {
	ID        uint           `gorm:"primarykey"` // 主键
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"uniqueIndex:uk_content_preset_external_id_deleted_at,priority:10"`        // 软删除时间（与 ExternalID 共同唯一）
	ExternalID *string        `gorm:"uniqueIndex:uk_content_preset_external_id_deleted_at,priority:1;size:64"` // 业务外部 ID（预设唯一标识）
	Name      string         `gorm:"size:128"`     // 预设名称（用于展示）
	RoleID    uint           `gorm:"index"`        // 关联岗位方向 ID
	IsActive  bool           `gorm:"default:true"` // 是否启用
}

func (ContentPreset) TableName() string {
	return "content_preset"
}
