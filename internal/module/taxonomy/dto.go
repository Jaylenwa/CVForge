package taxonomy

import "strings"

type JobCategoryDTO struct {
	ExternalID       string `json:"externalId"`
	Name             string `json:"name"`
	ParentExternalID string `json:"parentExternalId"`
	OrderNum         int    `json:"orderNum"`
	IsActive         bool   `json:"isActive"`
}

type JobRoleDTO struct {
	ExternalID         string   `json:"externalId"`
	CategoryExternalID string   `json:"categoryExternalId"`
	Name               string   `json:"name"`
	Tags               []string `json:"tags"`
	OrderNum           int      `json:"orderNum"`
	IsActive           bool     `json:"isActive"`
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

