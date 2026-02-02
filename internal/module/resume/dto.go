package resume

import (
	"sort"

	"openresume/internal/common"
)

type ResumeDTO struct {
	ID           uint              `json:"ID"`
	Title        string            `json:"Title"`
	TemplateID   string            `json:"TemplateID"`
	PresetID     *uint             `json:"PresetID"`
	RoleID       *uint             `json:"RoleID"`
	Language     string            `json:"Language"`
	LastModified int64             `json:"LastModified"`
	Personal     ResumePersonalDTO `json:"Personal"`
	Theme        ResumeThemeDTO    `json:"Theme"`
	Sections     []SectionDTO      `json:"Sections"`
}

type ResumePersonalDTO struct {
	FullName   string `json:"FullName"`
	Email      string `json:"Email"`
	Phone      string `json:"Phone"`
	AvatarURL  string `json:"AvatarURL"`
	Job        string `json:"Job"`
	City       string `json:"City"`
	Money      string `json:"Money"`
	JoinTime   string `json:"JoinTime"`
	Gender     string `json:"Gender"`
	Age        string `json:"Age"`
	Degree     string `json:"Degree"`
	CustomInfo string `json:"CustomInfo"`
}

type ResumeThemeDTO struct {
	Color    string `json:"Color"`
	Font     string `json:"Font"`
	Spacing  string `json:"Spacing"`
	FontSize string `json:"FontSize"`
}

type SectionDTO struct {
	ID        uint      `json:"ID"`
	Type      string    `json:"Type"`
	Title     string    `json:"Title"`
	IsVisible bool      `json:"IsVisible"`
	OrderNum  int       `json:"OrderNum"`
	Items     []ItemDTO `json:"Items"`
}

type ItemDTO struct {
	ID          uint   `json:"ID"`
	Title       string `json:"Title"`
	Subtitle    string `json:"Subtitle"`
	Major       string `json:"Major"`
	Degree      string `json:"Degree"`
	TimeStart   string `json:"TimeStart"`
	TimeEnd     string `json:"TimeEnd"`
	Today       bool   `json:"Today"`
	Description string `json:"Description"`
	OrderNum    int    `json:"OrderNum"`
}

func ToDTO(r Resume) ResumeDTO {
	dto := ResumeDTO{
		ID:           r.ID,
		Title:        r.Title,
		TemplateID:   r.TemplateID,
		PresetID:     r.PresetID,
		RoleID:       r.RoleID,
		Language:     r.Language,
		LastModified: r.LastModified,
		Personal: ResumePersonalDTO{
			FullName:   r.Personal.FullName,
			Email:      r.Personal.Email,
			Phone:      r.Personal.Phone,
			AvatarURL:  r.Personal.AvatarURL,
			Job:        r.Personal.Job,
			City:       r.Personal.City,
			Money:      r.Personal.Money,
			JoinTime:   r.Personal.JoinTime,
			Gender:     r.Personal.Gender,
			Age:        r.Personal.Age,
			Degree:     r.Personal.Degree,
			CustomInfo: r.Personal.CustomInfo,
		},
		Theme: ResumeThemeDTO{
			Color:    r.Theme.Color,
			Font:     r.Theme.Font,
			Spacing:  r.Theme.Spacing,
			FontSize: r.Theme.FontSize,
		},
	}
	secs := append([]ResumeSection(nil), r.Sections...)
	sort.SliceStable(secs, func(i, j int) bool { return secs[i].OrderNum < secs[j].OrderNum })
	for _, s := range secs {
		typ := s.Type
		if t, ok := common.NormalizeResumeSectionType(s.Type); ok {
			typ = string(t)
		}
		sec := SectionDTO{
			ID:        s.ID,
			Type:      typ,
			Title:     s.Title,
			IsVisible: s.IsVisible,
			OrderNum:  s.OrderNum,
		}
		items := append([]ResumeItem(nil), s.Items...)
		sort.SliceStable(items, func(i, j int) bool { return items[i].OrderNum < items[j].OrderNum })
		for _, it := range items {
			sec.Items = append(sec.Items, ItemDTO{
				ID:          it.ID,
				Title:       it.Title,
				Subtitle:    it.Subtitle,
				Major:       it.Major,
				Degree:      it.Degree,
				TimeStart:   it.TimeStart,
				TimeEnd:     it.TimeEnd,
				Today:       it.Today,
				Description: it.Description,
				OrderNum:    it.OrderNum,
			})
		}
		dto.Sections = append(dto.Sections, sec)
	}
	return dto
}

func ToPreviewDTO(r Resume, maxItems int) ResumeDTO {
	dto := ToDTO(r)
	if maxItems <= 0 {
		return dto
	}
	for i := range dto.Sections {
		if len(dto.Sections[i].Items) > maxItems {
			dto.Sections[i].Items = dto.Sections[i].Items[:maxItems]
		}
	}
	return dto
}
