package models

import (
	"time"

	"gorm.io/gorm"
)

type Template struct {
	ID         uint           `gorm:"primarykey"` // 主键
	CreatedAt  time.Time      // 创建时间
	UpdatedAt  time.Time      // 更新时间
	DeletedAt  gorm.DeletedAt `gorm:"uniqueIndex:uk_template_external_id_deleted_at,priority:10"`        // 软删除时间（与 ExternalID 共同唯一）
	ExternalID string         `gorm:"uniqueIndex:uk_template_external_id_deleted_at,priority:1;size:64"` // 业务外部 ID（模板唯一标识）
	Name       string         `gorm:"size:128"`                                                          // 模板名称
	UsageCount int            `gorm:"default:0"`                                                          // 使用次数（全局热度）
}

func (Template) TableName() string {
	return "template"
}
