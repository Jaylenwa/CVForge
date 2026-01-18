package models

import (
	"openresume/internal/common"

	"gorm.io/gorm"
)

type OAuthAccount struct {
	gorm.Model
	UserID          uint                `gorm:"index"`                                                                              // 关联用户 ID
	Provider        common.ProviderType `gorm:"size:32;index:idx_provider_openid,priority:1;index:idx_provider_unionid,priority:1"` // 登录提供方（如 wechat/github）
	ProviderOpenID  string              `gorm:"size:128;index:idx_provider_openid,priority:2"`                                      // 提供方 OpenID（同一 Provider 内唯一）
	ProviderUnionID string              `gorm:"size:128;index:idx_provider_unionid,priority:2"`                                     // 提供方 UnionID（同一 Provider 内可跨应用）
	RawProfileJSON  string              `gorm:"type:text"`                                                                          // 提供方原始用户信息（JSON 字符串）
}

func (OAuthAccount) TableName() string {
	return "oauth_account"
}
