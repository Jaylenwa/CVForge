package seed

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type ValidationError struct {
	Problems []string
}

func (e *ValidationError) Error() string {
	if e == nil || len(e.Problems) == 0 {
		return "seed validation failed"
	}
	var b strings.Builder
	b.WriteString("seed validation failed:")
	for _, p := range e.Problems {
		b.WriteString("\n- ")
		b.WriteString(p)
	}
	return b.String()
}

func (s SeedData) Validate() error {
	var problems []string

	categoryIDs := make(map[string]struct{}, len(s.Categories))
	for i, c := range s.Categories {
		if strings.TrimSpace(c.ExternalID) == "" {
			problems = append(problems, fmt.Sprintf("Categories[%d].ExternalID is empty", i))
			continue
		}
		if _, ok := categoryIDs[c.ExternalID]; ok {
			problems = append(problems, fmt.Sprintf("duplicate category ExternalID: %q", c.ExternalID))
			continue
		}
		categoryIDs[c.ExternalID] = struct{}{}
		if len(c.Names) == 0 {
			problems = append(problems, fmt.Sprintf("Categories[%d].Names is empty (%q)", i, c.ExternalID))
		}
	}
	for i, c := range s.Categories {
		if strings.TrimSpace(c.ParentExternalID) == "" {
			continue
		}
		if _, ok := categoryIDs[c.ParentExternalID]; !ok {
			problems = append(problems, fmt.Sprintf("Categories[%d] (%q) references missing ParentExternalID %q", i, c.ExternalID, c.ParentExternalID))
		}
	}

	roleIDs := make(map[string]struct{}, len(s.Roles))
	for i, r := range s.Roles {
		if strings.TrimSpace(r.ExternalID) == "" {
			problems = append(problems, fmt.Sprintf("Roles[%d].ExternalID is empty", i))
			continue
		}
		if _, ok := roleIDs[r.ExternalID]; ok {
			problems = append(problems, fmt.Sprintf("duplicate role ExternalID: %q", r.ExternalID))
			continue
		}
		roleIDs[r.ExternalID] = struct{}{}
		if strings.TrimSpace(r.CategoryExternalID) == "" {
			problems = append(problems, fmt.Sprintf("Roles[%d] (%q) CategoryExternalID is empty", i, r.ExternalID))
		} else if _, ok := categoryIDs[r.CategoryExternalID]; !ok {
			problems = append(problems, fmt.Sprintf("Roles[%d] (%q) references missing CategoryExternalID %q", i, r.ExternalID, r.CategoryExternalID))
		}
		if len(r.Names) == 0 {
			problems = append(problems, fmt.Sprintf("Roles[%d].Names is empty (%q)", i, r.ExternalID))
		}
	}

	presetIDs := make(map[string]struct{}, len(s.Presets))
	for i, p := range s.Presets {
		if strings.TrimSpace(p.ExternalID) == "" {
			problems = append(problems, fmt.Sprintf("Presets[%d].ExternalID is empty", i))
		} else if _, ok := presetIDs[p.ExternalID]; ok {
			problems = append(problems, fmt.Sprintf("duplicate preset ExternalID: %q", p.ExternalID))
		} else {
			presetIDs[p.ExternalID] = struct{}{}
		}
		if strings.TrimSpace(p.Name) == "" {
			problems = append(problems, fmt.Sprintf("Presets[%d] (%q) Name is empty", i, p.ExternalID))
		}
		if strings.TrimSpace(p.Language) == "" {
			problems = append(problems, fmt.Sprintf("Presets[%d] (%q) Language is empty", i, p.ExternalID))
		}
		if strings.TrimSpace(p.RoleCode) == "" {
			problems = append(problems, fmt.Sprintf("Presets[%d] (%q) RoleCode is empty", i, p.ExternalID))
		} else if _, ok := roleIDs[p.RoleCode]; !ok {
			problems = append(problems, fmt.Sprintf("Presets[%d] (%q) references missing RoleCode %q", i, p.ExternalID, p.RoleCode))
		}
		if strings.TrimSpace(p.DataJSON) == "" {
			problems = append(problems, fmt.Sprintf("Presets[%d] (%q) DataJSON is empty", i, p.ExternalID))
		} else if !json.Valid([]byte(p.DataJSON)) {
			problems = append(problems, fmt.Sprintf("Presets[%d] (%q) DataJSON is not valid json", i, p.ExternalID))
		}
	}

	if len(problems) > 0 {
		sort.Strings(problems)
		return &ValidationError{Problems: problems}
	}
	return nil
}

