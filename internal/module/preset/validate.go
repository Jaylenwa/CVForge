package preset

import (
	"encoding/json"
	"errors"
	"strings"

	"cvforge/internal/common"
)

func ValidatePresetDataJSON(dataJSON string) error {
	if strings.TrimSpace(dataJSON) == "" {
		return nil
	}
	if !json.Valid([]byte(dataJSON)) {
		return errors.New("invalid_json")
	}
	var root any
	if err := json.Unmarshal([]byte(dataJSON), &root); err != nil {
		return errors.New("invalid_json")
	}
	m, ok := root.(map[string]any)
	if !ok {
		return errors.New("invalid_schema")
	}
	if v, ok := m["sections"]; ok && v != nil {
		secs, ok := v.([]any)
		if !ok {
			return errors.New("invalid_schema_sections")
		}
		for _, s := range secs {
			sm, ok := s.(map[string]any)
			if !ok {
				return errors.New("invalid_schema_sections")
			}
			if id, ok := sm["id"]; ok && id != nil {
				if _, ok := id.(string); !ok {
					return errors.New("invalid_schema_section_id")
				}
			}
			if typ, ok := sm["type"]; ok && typ != nil {
				str, ok := typ.(string)
				if !ok {
					return errors.New("invalid_schema_section_type")
				}
				if _, ok := common.NormalizeResumeSectionType(str); !ok {
					return errors.New("invalid_schema_section_type_value")
				}
			}
			if items, ok := sm["items"]; ok && items != nil {
				if _, ok := items.([]any); !ok {
					return errors.New("invalid_schema_section_items")
				}
			}
		}
	}
	return nil
}
