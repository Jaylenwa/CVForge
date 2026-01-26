package models

import (
	"time"

	"gorm.io/gorm"
)

type TemplateVariant struct {
	ID                       uint           `gorm:"primarykey"` // 主键
	CreatedAt                time.Time      // 创建时间
	UpdatedAt                time.Time      // 更新时间
	DeletedAt                gorm.DeletedAt `gorm:"uniqueIndex:uk_template_variant_role_preset_tpl_deleted_at,priority:10"` // 软删除时间
	Name                     string         `gorm:"size:128"`                                                                  // 变体名称（用于展示）
	LayoutTemplateExternalID string         `gorm:"size:64;index;uniqueIndex:uk_template_variant_role_preset_tpl_deleted_at,priority:3"` // 关联布局模板 ExternalID（Template.ExternalID）
	PresetID                 uint           `gorm:"index;uniqueIndex:uk_template_variant_role_preset_tpl_deleted_at,priority:2"`          // 关联内容预设 ID
	RoleID                   uint           `gorm:"index;uniqueIndex:uk_template_variant_role_preset_tpl_deleted_at,priority:1"`          // 关联岗位方向 ID
	Tags                     string         `gorm:"size:256"`                                                                  // 标签（逗号分隔）
	UsageCount               int            `gorm:"default:0"`                                                                 // 使用次数（热度）
	IsPremium                bool           `gorm:"default:false"`                                                             // 是否付费变体
	IsActive                 bool           `gorm:"default:true"`                                                              // 是否启用
}

func (TemplateVariant) TableName() string {
	return "template_variant"
}
