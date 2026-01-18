package catalog

type AdminPageResp[T any] struct {
	Items    []T   `json:"items"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

type GenerateMode string

const (
	GenerateModeSkip   GenerateMode = "skip"
	GenerateModeUpdate GenerateMode = "update"
)

type GenerateVariantsReq struct {
	RoleID            string   `json:"roleId"`
	PresetID          string   `json:"presetId"`
	LayoutTemplateIDs []string `json:"layoutTemplateIds"`
	NamePrefix        string   `json:"namePrefix"`
	Tags              string   `json:"tags"`
	IsPremium         *bool    `json:"isPremium"`
	IsActive          *bool    `json:"isActive"`
	Mode              string   `json:"mode"`
}

type GenerateVariantsResult struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
	Skipped int `json:"skipped"`
	Failed  int `json:"failed"`
	Items   []struct {
		LayoutTemplateID string `json:"layoutTemplateId"`
		ExternalID       string `json:"externalId"`
		Action           string `json:"action"`
		Error            string `json:"error,omitempty"`
	} `json:"items"`
}

