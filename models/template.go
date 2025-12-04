package models

import "gorm.io/gorm"

type Template struct {
    gorm.Model
    ExternalID string `gorm:"uniqueIndex;size:64"`
    Name       string `gorm:"size:128"`
    Thumbnail  string `gorm:"size:512"`
    Tags       string `gorm:"size:256"` // comma-separated for简单实现
    Popularity int
    IsPremium  bool
    Category   string `gorm:"size:64"`
    Level      string `gorm:"size:64"`
}
