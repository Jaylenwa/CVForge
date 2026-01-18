package models

import (
	"time"

	"gorm.io/gorm"
)

type JobRole struct {
	ID                 uint           `gorm:"primarykey"` // 主键
	CreatedAt          time.Time      // 创建时间
	UpdatedAt          time.Time      // 更新时间
	DeletedAt          gorm.DeletedAt `gorm:"uniqueIndex:uk_job_role_external_id_deleted_at,priority:10"`        // 软删除时间（与 ExternalID 共同唯一）
	ExternalID         string         `gorm:"uniqueIndex:uk_job_role_external_id_deleted_at,priority:1;size:64"` // 业务外部 ID（岗位方向唯一标识）
	CategoryExternalID string         `gorm:"size:64;index"`                                                     // 所属岗位分类 ExternalID
	Name               string         `gorm:"size:128"`                                                          // 岗位方向名称（如 Java 开发工程师）
	Tags               string         `gorm:"size:256"`                                                          // 标签（逗号分隔）
	OrderNum           int            `gorm:"default:0"`                                                         // 排序号（越小越靠前）
	IsActive           bool           `gorm:"default:true"`                                                      // 是否启用
}

func (JobRole) TableName() string {
	return "job_role"
}
