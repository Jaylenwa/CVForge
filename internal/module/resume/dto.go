package resume

type ResumeDTO struct {
	ExternalID   string          `json:"ExternalID"`
	Title        string          `json:"Title"`
	TemplateID   string          `json:"TemplateID"`
	LastModified int64           `json:"LastModified"`
	ThemeColor   string          `json:"ThemeColor"`
	ThemeFont    string          `json:"ThemeFont"`
	ThemeSpacing string          `json:"ThemeSpacing"`
	FullName     string          `json:"FullName"`
	Email        string          `json:"Email"`
	Phone        string          `json:"Phone"`
	AvatarURL    string          `json:"AvatarURL"`
	JobTitle     string          `json:"JobTitle"`
	Gender       string          `json:"Gender"`
	Age          string          `json:"Age"`
	MaritalStatus   string       `json:"MaritalStatus"`
	PoliticalStatus string       `json:"PoliticalStatus"`
	Birthplace      string       `json:"Birthplace"`
	Ethnicity       string       `json:"Ethnicity"`
	Height          string       `json:"Height"`
	Weight          string       `json:"Weight"`
	CustomInfo      string       `json:"CustomInfo"`
	Sections     []ResumeSection `json:"Sections"`
}

func ToDTO(r Resume) ResumeDTO {
	return ResumeDTO{
		ExternalID:   r.ExternalID,
		Title:        r.Title,
		TemplateID:   r.TemplateID,
		LastModified: r.LastModified,
		ThemeColor:   r.Theme.Color,
		ThemeFont:    r.Theme.Font,
		ThemeSpacing: r.Theme.Spacing,
		FullName:     r.Personal.FullName,
		Email:        r.Personal.Email,
		Phone:        r.Personal.Phone,
		AvatarURL:    r.Personal.AvatarURL,
		JobTitle:     r.Personal.JobTitle,
		Gender:       r.Personal.Gender,
		Age:          r.Personal.Age,
		MaritalStatus:   r.Personal.MaritalStatus,
		PoliticalStatus: r.Personal.PoliticalStatus,
		Birthplace:      r.Personal.Birthplace,
		Ethnicity:       r.Personal.Ethnicity,
		Height:          r.Personal.Height,
		Weight:          r.Personal.Weight,
		CustomInfo:      r.Personal.CustomInfo,
		Sections:     r.Sections,
	}
}
