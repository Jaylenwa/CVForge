package models

import "gorm.io/gorm"

type Config struct {
	gorm.Model
	ConfigKey   string `gorm:"column:config_key;uniqueIndex;size:100" json:"key"` // 配置键（唯一）
	ConfigValue string `gorm:"column:config_value;type:text" json:"value"`        // 配置值（字符串保存）
	Description string `gorm:"size:255" json:"description"`                       // 配置说明
	Type        string `gorm:"size:20;default:string" json:"type"`                // 配置类型（string/bool/int/json）
}

func (Config) TableName() string {
	return "config"
}
