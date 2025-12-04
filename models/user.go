package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Email       string `gorm:"uniqueIndex;size:191"`
    PasswordHash string `gorm:"size:255"`
    Name        string `gorm:"size:128"`
    AvatarURL   string `gorm:"size:512"`
    Language    string `gorm:"size:8"`
}
