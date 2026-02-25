package models

import (
	"openresume/internal/common"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        *string     `gorm:"uniqueIndex;size:191;null"` // 登录邮箱（唯一，可为空）
	PasswordHash string      `gorm:"size:255"`                  // 密码哈希（不可逆）
	Name         string      `gorm:"size:128"`                  // 昵称/姓名
	AvatarURL    string      `gorm:"size:512"`                  // 头像 URL
	Language     string      `gorm:"size:8"`                    // 语言偏好（zh/en）
	Role         common.Role `gorm:"size:16;default:user"`      // 角色（user/admin）
	IsActive     bool        `gorm:"default:true"`              // 是否启用（封禁时为 false）
	LastLoginAt  *time.Time  // 最近登录时间
	TokenVersion int         `gorm:"default:1"` // token 版本（用于全局登出与撤销）
}

func (User) TableName() string {
	return "user"
}
