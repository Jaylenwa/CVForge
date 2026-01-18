package models

import "gorm.io/gorm"

type Resume struct {
	gorm.Model
	ExternalID   string          `gorm:"uniqueIndex;size:64"` // 业务外部 ID（用于 URL/接口）
	UserID       uint            // 所属用户 ID
	Title        string          `gorm:"size:191"`      // 简历标题
	TemplateID   string          `gorm:"size:64"`       // 布局模板 ExternalID
	VariantID    string          `gorm:"size:64;index"` // 模板变体 ExternalID（岗位维度）
	PresetID     string          `gorm:"size:64;index"` // 内容预设 ExternalID
	RoleID       string          `gorm:"size:64;index"` // 岗位方向 ExternalID
	Language     string          `gorm:"size:8"`        // 语言（zh/en）
	LastModified int64           // 最近更新时间戳（毫秒）
	Personal     ResumePersonal  `gorm:"constraint:OnDelete:CASCADE;foreignKey:ResumeID"` // 个人信息（1:1）
	Theme        ResumeTheme     `gorm:"constraint:OnDelete:CASCADE;foreignKey:ResumeID"` // 主题配置（1:1）
	Sections     []ResumeSection `gorm:"constraint:OnDelete:CASCADE;foreignKey:ResumeID"` // 简历模块列表（1:N）
}

func (Resume) TableName() string {
	return "resume"
}

type ResumePersonal struct {
	gorm.Model
	ResumeID   uint   `gorm:"uniqueIndex"` // 所属简历 ID（1:1）
	FullName   string `gorm:"size:128"`    // 姓名
	Email      string `gorm:"size:191"`    // 邮箱
	Phone      string `gorm:"size:64"`     // 手机号
	AvatarURL  string `gorm:"size:512"`    // 头像 URL
	Job        string `gorm:"size:128"`    // 求职岗位/当前岗位
	City       string `gorm:"size:128"`    // 城市
	Money      string `gorm:"size:64"`     // 期望薪资
	JoinTime   string `gorm:"size:64"`     // 到岗时间
	Gender     string `gorm:"size:32"`     // 性别
	Age        string `gorm:"size:32"`     // 年龄
	Degree     string `gorm:"size:64"`     // 学历
	CustomInfo string `gorm:"type:text"`   // 自定义信息（JSON 字符串）
}

func (ResumePersonal) TableName() string {
	return "resume_personal"
}

type ResumeTheme struct {
	gorm.Model
	ResumeID uint   `gorm:"uniqueIndex"` // 所属简历 ID（1:1）
	Color    string `gorm:"size:32"`     // 主色（HEX）
	Font     string `gorm:"size:64"`     // 字体（内部标识）
	Spacing  string `gorm:"size:32"`     // 行距/间距配置
	FontSize string `gorm:"size:16"`     // 字号配置
}

func (ResumeTheme) TableName() string {
	return "resume_theme"
}

type ResumeSection struct {
	gorm.Model
	ResumeID   uint         // 所属简历 ID
	ExternalID string       `gorm:"size:64"`  // 模块外部 ID（前端 section.id）
	Type       string       `gorm:"size:32"`  // 模块类型（如 Experience/Education/Skills）
	Title      string       `gorm:"size:128"` // 模块标题
	IsVisible  bool         // 是否显示
	OrderNum   int          // 排序号（越小越靠前）
	Items      []ResumeItem `gorm:"constraint:OnDelete:CASCADE;foreignKey:SectionID"` // 模块条目（1:N）
}

func (ResumeSection) TableName() string {
	return "resume_section"
}

type ResumeItem struct {
	gorm.Model
	SectionID   uint   // 所属模块 ID
	ExternalID  string `gorm:"size:64"`  // 条目外部 ID（前端 item.id）
	Title       string `gorm:"size:128"` // 主标题（公司/学校/项目等）
	Subtitle    string `gorm:"size:128"` // 副标题（岗位/专业等）
	Major       string `gorm:"size:128"` // 专业（教育模块常用）
	Degree      string `gorm:"size:128"` // 学位/学历（教育模块常用）
	TimeStart   string `gorm:"size:7"`   // 开始时间（YYYY-MM）
	TimeEnd     string `gorm:"size:7"`   // 结束时间（YYYY-MM）
	Today       bool   // 是否至今
	Description string `gorm:"type:text"` // 详情描述（支持富文本）
	OrderNum    int    // 排序号（越小越靠前）
}

func (ResumeItem) TableName() string {
	return "resume_item"
}
