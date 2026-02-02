package library

type TemplateLibraryItemDTO struct {
	TemplateExternalID string `json:"templateExternalId"`
	Name               string `json:"name"`
	UsageCount         int    `json:"usageCount"`
	GlobalUsageCount   int    `json:"globalUsageCount"`
	PresetID           uint   `json:"presetId,omitempty"`
	RoleID             uint   `json:"roleId,omitempty"`
}
