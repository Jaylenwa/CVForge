package library

import (
	"strings"

	"gorm.io/gorm"
)

type TemplateLibraryItem struct {
	TemplateExternalID string
	Name               string
	UsageCount         int
	GlobalUsageCount   int
	PresetID           uint
	RoleID             uint
	IsPremium          bool
}

func normalizeTemplateListSort(raw string) string {
	v := strings.TrimSpace(strings.ToLower(raw))
	switch v {
	case "hot", "new", "name":
		return v
	default:
		return "hot"
	}
}

func normalizeLanguage(raw string) string {
	v := strings.TrimSpace(strings.ToLower(raw))
	if v == "en" {
		return "en"
	}
	if v == "zh" {
		return "zh"
	}
	return "zh"
}

func (r *Repo) pickPresetID(roleID uint, language string) uint {
	language = normalizeLanguage(language)
	var presetID uint

	q := r.db.Table("content_preset cp").
		Select("cp.id").
		Joins("JOIN content_preset_i18n cpi ON cpi.content_preset_id = cp.id AND cpi.deleted_at IS NULL").
		Where("cp.is_active = ?", true).
		Where("cp.deleted_at IS NULL")
	if roleID != 0 {
		q = q.Where("cp.role_id = ?", roleID)
	}

	_ = q.Where("cpi.language = ?", language).
		Order("cp.id asc").
		Limit(1).
		Scan(&presetID).Error
	if presetID != 0 {
		return presetID
	}

	_ = q.Where("cpi.language = ?", "zh").
		Order("cp.id asc").
		Limit(1).
		Scan(&presetID).Error
	if presetID != 0 {
		return presetID
	}

	_ = q.Order("cp.id asc").Limit(1).Scan(&presetID).Error
	return presetID
}

func (r *Repo) ListTemplateLibraryItems(roleID uint, language string, sort string) ([]TemplateLibraryItem, error) {
	sort = normalizeTemplateListSort(sort)
	presetID := r.pickPresetID(roleID, language)
	if roleID != 0 && presetID == 0 {
		return []TemplateLibraryItem{}, nil
	}
	if roleID == 0 && presetID == 0 {
		return []TemplateLibraryItem{}, nil
	}

	type row struct {
		TemplateExternalID string `gorm:"column:template_external_id"`
		Name               string `gorm:"column:name"`
		UsageCount         int    `gorm:"column:usage_count"`
		GlobalUsageCount   int    `gorm:"column:global_usage_count"`
	}
	var rows []row
	db := r.db.Table("template t").
		Where("t.deleted_at IS NULL")

	if roleID != 0 {
		db = db.Joins("LEFT JOIN role_template_usage rtu ON rtu.template_external_id = t.external_id AND rtu.role_id = ?", roleID).
			Select("t.external_id as template_external_id, t.name, COALESCE(rtu.usage_count, 0) as usage_count, t.usage_count as global_usage_count")
	} else {
		db = db.Select("t.external_id as template_external_id, t.name, t.usage_count as usage_count, t.usage_count as global_usage_count")
	}

	switch sort {
	case "new":
		db = db.Order("t.id desc")
	case "name":
		db = db.Order("t.name asc").Order("t.id asc")
	default:
		if roleID != 0 {
			db = db.Order("COALESCE(rtu.usage_count, 0) desc").Order("t.id desc")
		} else {
			db = db.Order("t.usage_count desc").Order("t.id desc")
		}
	}

	if err := db.Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]TemplateLibraryItem, 0, len(rows))
	for _, it := range rows {
		out = append(out, TemplateLibraryItem{
			TemplateExternalID: it.TemplateExternalID,
			Name:               it.Name,
			UsageCount:         it.UsageCount,
			GlobalUsageCount:   it.GlobalUsageCount,
			PresetID:           presetID,
			RoleID:             roleID,
			IsPremium:          false,
		})
	}
	return out, nil
}

func (r *Repo) PatchTemplateUsageTx(tx *gorm.DB, roleID uint, templateExternalID string, delta int) error {
	if tx == nil || roleID == 0 || templateExternalID == "" || delta == 0 {
		return nil
	}
	return tx.Exec(`UPDATE role_template_usage SET usage_count = usage_count + ? WHERE role_id = ? AND template_external_id = ?`, delta, roleID, templateExternalID).Error
}
