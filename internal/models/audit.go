package models

import "gorm.io/gorm"

type AuditLog struct {
	gorm.Model
	ActorID    uint   // 操作人用户 ID
	Action     string `gorm:"size:64"`   // 操作动作（如 admin.user.ban）
	TargetType string `gorm:"size:32"`   // 目标类型（如 user/resume/template）
	TargetID   string `gorm:"size:128"`  // 目标 ID（通常为 external_id 或数字 ID）
	Metadata   string `gorm:"type:text"` // 扩展信息（JSON 字符串）
	IP         string `gorm:"size:64"`   // 请求 IP
	UA         string `gorm:"size:256"`  // User-Agent
}

func (AuditLog) TableName() string {
	return "audit_log"
}
