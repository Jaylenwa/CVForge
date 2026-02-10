package models

import (
	"time"

	"gorm.io/gorm"
)

type JobCategory struct {
	ID         uint           `gorm:"primarykey"` // 主键
	CreatedAt  time.Time      // 创建时间
	UpdatedAt  time.Time      // 更新时间
	DeletedAt  gorm.DeletedAt `gorm:"uniqueIndex:uk_job_category_external_id_deleted_at,priority:10"`        // 软删除时间（与 ExternalID 共同唯一）
	ExternalID *string        `gorm:"uniqueIndex:uk_job_category_external_id_deleted_at,priority:1;size:64"` // 业务外部 ID（分类唯一标识）
	ParentID   *uint          `gorm:"index"`                                                                 // 父级分类 ID（用于多级分类）
	OrderNum   int            `gorm:"default:0"`                                                             // 排序号（越小越靠前）
	IsActive   bool           `gorm:"default:true"`                                                          // 是否启用
}

func (JobCategory) TableName() string {
	return "job_category"
}
