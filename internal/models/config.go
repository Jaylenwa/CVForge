package models

import "gorm.io/gorm"

type Config struct {
	gorm.Model
	ConfigKey   string `gorm:"column:config_key;uniqueIndex;size:100" json:"key"`
	ConfigValue string `gorm:"column:config_value;type:text" json:"value"`
	Description string `gorm:"size:255" json:"description"`
	Type        string `gorm:"size:20;default:string" json:"type"` // string, bool, int, json
}

func (Config) TableName() string {
	return "config"
}
