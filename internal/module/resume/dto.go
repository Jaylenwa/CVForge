package resume

type ResumeDTO struct {
	ExternalID      string       `json:"ExternalID"`
	Title           string       `json:"Title"`
	TemplateID      string       `json:"TemplateID"`
	LastModified    int64        `json:"LastModified"`
	ThemeColor      string       `json:"ThemeColor"`
	ThemeFont       string       `json:"ThemeFont"`
	ThemeSpacing    string       `json:"ThemeSpacing"`
	FullName        string       `json:"FullName"`
	Email           string       `json:"Email"`
	Phone           string       `json:"Phone"`
	AvatarURL       string       `json:"AvatarURL"`
	JobTitle        string       `json:"JobTitle"`
	Gender          string       `json:"Gender"`
	Age             string       `json:"Age"`
	MaritalStatus   string       `json:"MaritalStatus"`
	PoliticalStatus string       `json:"PoliticalStatus"`
	Birthplace      string       `json:"Birthplace"`
	Ethnicity       string       `json:"Ethnicity"`
	Height          string       `json:"Height"`
	Weight          string       `json:"Weight"`
	CustomInfo      string       `json:"CustomInfo"`
	Sections        []SectionDTO `json:"Sections"`
}

type SectionDTO struct {
	ExternalID string    `json:"ExternalID"`
	Type       string    `json:"Type"`
	Title      string    `json:"Title"`
	IsVisible  bool      `json:"IsVisible"`
	Items      []ItemDTO `json:"Items"`
}

type ItemDTO struct {
	ExternalID  string `json:"ExternalID"`
	Title       string `json:"Title"`
	Subtitle    string `json:"Subtitle"`
	TimeStart   string `json:"TimeStart"`
	TimeEnd     string `json:"TimeEnd"`
	Today       bool   `json:"Today"`
	Description string `json:"Description"`
}

func ToDTO(r Resume) ResumeDTO {
	dto := ResumeDTO{
		ExternalID:      r.ExternalID,
		Title:           r.Title,
		TemplateID:      r.TemplateID,
		LastModified:    r.LastModified,
		ThemeColor:      r.Theme.Color,
		ThemeFont:       r.Theme.Font,
		ThemeSpacing:    r.Theme.Spacing,
		FullName:        r.Personal.FullName,
		Email:           r.Personal.Email,
		Phone:           r.Personal.Phone,
		AvatarURL:       r.Personal.AvatarURL,
		JobTitle:        r.Personal.JobTitle,
		Gender:          r.Personal.Gender,
		Age:             r.Personal.Age,
		MaritalStatus:   r.Personal.MaritalStatus,
		PoliticalStatus: r.Personal.PoliticalStatus,
		Birthplace:      r.Personal.Birthplace,
		Ethnicity:       r.Personal.Ethnicity,
		Height:          r.Personal.Height,
		Weight:          r.Personal.Weight,
		CustomInfo:      r.Personal.CustomInfo,
	}
	for _, s := range r.Sections {
		sec := SectionDTO{
			ExternalID: s.ExternalID,
			Type:       s.Type,
			Title:      s.Title,
			IsVisible:  s.IsVisible,
		}
		for _, it := range s.Items {
			sec.Items = append(sec.Items, ItemDTO{
				ExternalID:  it.ExternalID,
				Title:       it.Title,
				Subtitle:    it.Subtitle,
				TimeStart:   it.TimeStart,
				TimeEnd:     it.TimeEnd,
				Today:       it.Today,
				Description: it.Description,
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
