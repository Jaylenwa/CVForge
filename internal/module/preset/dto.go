package preset

import "strings"

type ContentPresetDTO struct {
	ExternalID     string   `json:"externalId"`
	Name           string   `json:"name"`
	Language       string   `json:"language"`
	RoleExternalID string   `json:"roleExternalId"`
	Tags           []string `json:"tags"`
	DataJSON       string   `json:"dataJson"`
	IsActive       bool     `json:"isActive"`
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

