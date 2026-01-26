package library

import "strings"

type TemplateVariantDTO struct {
	ID                       uint     `json:"id"`
	Name                     string   `json:"name"`
	LayoutTemplateExternalID string   `json:"layoutTemplateExternalId"`
	PresetID                 uint     `json:"presetId"`
	RoleID                   uint     `json:"roleId"`
	Tags                     []string `json:"tags"`
	UsageCount               int      `json:"usageCount"`
	IsPremium                bool     `json:"isPremium"`
	IsActive                 bool     `json:"isActive"`
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
