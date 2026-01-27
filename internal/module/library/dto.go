package library

import "strings"

type TemplateLibraryItemDTO struct {
	TemplateExternalID string   `json:"templateExternalId"`
	Name               string   `json:"name"`
	Tags               []string `json:"tags"`
	UsageCount         int      `json:"usageCount"`
	GlobalUsageCount   int      `json:"globalUsageCount"`
	PresetID           uint     `json:"presetId,omitempty"`
	RoleID             uint     `json:"roleId,omitempty"`
	IsPremium          bool     `json:"isPremium"`
}

func splitTags(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
