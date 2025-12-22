package models

import "gorm.io/gorm"

type Config struct {
	gorm.Model
	Key         string `gorm:"uniqueIndex;size:100" json:"key"`
	Value       string `gorm:"type:text" json:"value"`
	Description string `gorm:"size:255" json:"description"`
	Type        string `gorm:"size:20;default:string" json:"type"` // string, bool, int, json
}

func (Config) TableName() string {
	return "config"
}
