package models

import (
	"time"

	"gorm.io/gorm"
)

type JobRole struct {
	ID         uint           `gorm:"primarykey"` // 主键
	CreatedAt  time.Time      // 创建时间
	UpdatedAt  time.Time      // 更新时间
	DeletedAt  gorm.DeletedAt `gorm:"index"`        // 软删除时间
	CategoryID uint           `gorm:"index"`        // 所属岗位分类 ID
	OrderNum   int            `gorm:"default:0"`    // 排序号（越小越靠前）
	IsActive   bool           `gorm:"default:true"` // 是否启用
}

func (JobRole) TableName() string {
	return "job_role"
}
