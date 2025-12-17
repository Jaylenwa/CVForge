package models

import "gorm.io/gorm"

type Resume struct {
	gorm.Model
	ExternalID      string `gorm:"uniqueIndex;size:64"`
	UserID          uint
	Title           string `gorm:"size:191"`
	TemplateID      string `gorm:"size:64"`
	ThemeColor      string `gorm:"size:32"`
	ThemeFont       string `gorm:"size:64"`
	ThemeSpacing    string `gorm:"size:32"`
	LastModified    int64
	FullName        string          `gorm:"size:128"`
	Email           string          `gorm:"size:191"`
	Phone           string          `gorm:"size:64"`
	AvatarURL       string          `gorm:"size:512"`
	JobTitle        string          `gorm:"size:128"`
	Gender          string          `gorm:"size:32"`
	Age             string          `gorm:"size:32"`
	MaritalStatus   string          `gorm:"size:64"`
	PoliticalStatus string          `gorm:"size:64"`
	Birthplace      string          `gorm:"size:128"`
	Ethnicity       string          `gorm:"size:64"`
	Height          string          `gorm:"size:32"`
	Weight          string          `gorm:"size:32"`
	CustomInfo      string          `gorm:"type:text"`
	Sections        []ResumeSection `gorm:"constraint:OnDelete:CASCADE;foreignKey:ResumeID"`
}

type ResumeSection struct {
	gorm.Model
	ResumeID   uint
	ExternalID string `gorm:"size:64"`
	Type       string `gorm:"size:32"`
	Title      string `gorm:"size:128"`
	IsVisible  bool
	OrderNum   int
	Items      []ResumeItem `gorm:"constraint:OnDelete:CASCADE;foreignKey:SectionID"`
}

type ResumeItem struct {
	gorm.Model
	SectionID   uint
	ExternalID  string `gorm:"size:64"`
	Title       string `gorm:"size:128"`
	Subtitle    string `gorm:"size:128"`
	DateRange   string `gorm:"size:64"`
	Description string `gorm:"type:text"`
	OrderNum    int
}
