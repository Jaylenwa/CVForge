package models

import "gorm.io/gorm"

type OAuthAccount struct {
	gorm.Model
	UserID           uint   `gorm:"index"`
	Provider         string `gorm:"size:32;index:idx_provider_openid,priority:1;index:idx_provider_unionid,priority:1"`
	ProviderOpenID   string `gorm:"size:128;index:idx_provider_openid,priority:2"`
	ProviderUnionID  string `gorm:"size:128;index:idx_provider_unionid,priority:2"`
	RawProfileJSON   string `gorm:"type:text"`
}
