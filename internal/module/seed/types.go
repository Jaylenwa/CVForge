package seed

type SeedData struct {
	Categories []SeedJobCategory `json:"Categories"`
	Roles      []SeedJobRole     `json:"Roles"`
	Presets    []SeedContentPreset `json:"Presets"`
}

type SeedJobCategory struct {
	ExternalID       string            `json:"ExternalID"`
	Names            map[string]string `json:"Names"`
	ParentExternalID string            `json:"ParentExternalID"`
	OrderNum         int               `json:"OrderNum"`
	IsActive         bool              `json:"IsActive"`
}

type SeedJobRole struct {
	ExternalID         string            `json:"ExternalID"`
	CategoryExternalID string            `json:"CategoryExternalID"`
	Names              map[string]string `json:"Names"`
	OrderNum           int               `json:"OrderNum"`
	IsActive           bool              `json:"IsActive"`
}

type SeedContentPreset struct {
	ExternalID string `json:"ExternalID"`
	Name       string `json:"Name"`
	Language   string `json:"Language"`
	RoleCode   string `json:"RoleCode"`
	DataJSON   string `json:"DataJSON"`
	IsActive   bool   `json:"IsActive"`
}

type ImportCounts struct {
	JobCategories  int `json:"jobCategories"`
	JobRoles       int `json:"jobRoles"`
	ContentPresets int `json:"contentPresets"`
}

