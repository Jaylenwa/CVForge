package models

import (
	"time"

	"gorm.io/gorm"
)

type TemplateVariant struct {
	ID                       uint           `gorm:"primarykey"` // 主键
	CreatedAt                time.Time      // 创建时间
	UpdatedAt                time.Time      // 更新时间
	DeletedAt                gorm.DeletedAt `gorm:"uniqueIndex:uk_template_variant_external_id_deleted_at,priority:10"`        // 软删除时间（与 ExternalID 共同唯一）
	ExternalID               string         `gorm:"uniqueIndex:uk_template_variant_external_id_deleted_at,priority:1;size:64"` // 业务外部 ID（变体唯一标识）
	Name                     string         `gorm:"size:128"`                                                                  // 变体名称（用于展示）
	LayoutTemplateExternalID string         `gorm:"size:64;index"`                                                             // 关联布局模板 ExternalID（Template.ExternalID）
	PresetExternalID         string         `gorm:"size:64;index"`                                                             // 关联内容预设 ExternalID（ContentPreset.ExternalID）
	RoleExternalID           string         `gorm:"size:64;index"`                                                             // 关联岗位方向 ExternalID（JobRole.ExternalID）
	Tags                     string         `gorm:"size:256"`                                                                  // 标签（逗号分隔）
	UsageCount               int            `gorm:"default:0"`                                                                 // 使用次数（热度）
	IsPremium                bool           `gorm:"default:false"`                                                             // 是否付费变体
	IsActive                 bool           `gorm:"default:true"`                                                              // 是否启用
}

func (TemplateVariant) TableName() string {
	return "template_variant"
}
