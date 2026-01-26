package taxonomy

import "strings"

type JobCategoryDTO struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	ParentID         *uint  `json:"parentId"`
	OrderNum         int    `json:"orderNum"`
	IsActive         bool   `json:"isActive"`
}

type JobRoleDTO struct {
	ID         uint   `json:"id"`
	CategoryID uint   `json:"categoryId"`
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
