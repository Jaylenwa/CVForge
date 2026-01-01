package models

import "gorm.io/gorm"

type Resume struct {
	gorm.Model
	ExternalID   string `gorm:"uniqueIndex;size:64"`
	UserID       uint
	Title        string `gorm:"size:191"`
	TemplateID   string `gorm:"size:64"`
	LastModified int64
	Personal     ResumePersonal  `gorm:"constraint:OnDelete:CASCADE;foreignKey:ResumeID"`
	Theme        ResumeTheme     `gorm:"constraint:OnDelete:CASCADE;foreignKey:ResumeID"`
	Sections     []ResumeSection `gorm:"constraint:OnDelete:CASCADE;foreignKey:ResumeID"`
}

func (Resume) TableName() string {
	return "resume"
}

type ResumePersonal struct {
	gorm.Model
	ResumeID        uint   `gorm:"uniqueIndex"`
	FullName        string `gorm:"size:128"`
	Email           string `gorm:"size:191"`
	Phone           string `gorm:"size:64"`
	AvatarURL       string `gorm:"size:512"`
	JobTitle        string `gorm:"size:128"`
	Gender          string `gorm:"size:32"`
	Age             string `gorm:"size:32"`
	MaritalStatus   string `gorm:"size:64"`
	PoliticalStatus string `gorm:"size:64"`
	Birthplace      string `gorm:"size:128"`
	Ethnicity       string `gorm:"size:64"`
	Height          string `gorm:"size:32"`
	Weight          string `gorm:"size:32"`
	CustomInfo      string `gorm:"type:text"`
}

func (ResumePersonal) TableName() string {
	return "resume_personal"
}

type ResumeTheme struct {
	gorm.Model
	ResumeID uint   `gorm:"uniqueIndex"`
	Color    string `gorm:"size:32"`
	Font     string `gorm:"size:64"`
	Spacing  string `gorm:"size:32"`
}

func (ResumeTheme) TableName() string {
	return "resume_theme"
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

func (ResumeSection) TableName() string {
	return "resume_section"
}

type ResumeItem struct {
	gorm.Model
	SectionID   uint
	ExternalID  string `gorm:"size:64"`
	Title       string `gorm:"size:128"`
	Subtitle    string `gorm:"size:128"`
	TimeStart   string `gorm:"size:7"`
	TimeEnd     string `gorm:"size:7"`
	Today       bool
	Description string `gorm:"type:text"`
	OrderNum    int
}

func (ResumeItem) TableName() string {
	return "resume_item"
}
