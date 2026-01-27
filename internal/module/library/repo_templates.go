package library

import "gorm.io/gorm"

type TemplateLibraryItem struct {
	TemplateExternalID string
	Name               string
	Tags               string
	UsageCount         int
	GlobalUsageCount   int
	PresetID           uint
	RoleID             uint
	IsPremium          bool
}

func (r *Repo) ListTemplateLibraryItems(roleID uint) ([]TemplateLibraryItem, error) {
	var presetID uint
	if roleID != 0 {
		_ = r.db.Table("content_preset").
			Select("id").
			Where("role_id = ?", roleID).
			Where("is_active = ?", true).
			Where("deleted_at IS NULL").
			Order("id asc").
			Limit(1).
			Scan(&presetID).Error
	}

	type row struct {
		TemplateExternalID string `gorm:"column:template_external_id"`
		Name               string `gorm:"column:name"`
		Tags               string `gorm:"column:tags"`
		UsageCount         int    `gorm:"column:usage_count"`
		GlobalUsageCount   int    `gorm:"column:global_usage_count"`
	}
	var rows []row
	db := r.db.Table("template t").
		Where("t.deleted_at IS NULL")

	if roleID != 0 {
		db = db.Joins("LEFT JOIN role_template_usage rtu ON rtu.template_external_id = t.external_id AND rtu.role_id = ?", roleID).
			Select("t.external_id as template_external_id, t.name, t.tags, COALESCE(rtu.usage_count, 0) as usage_count, t.usage_count as global_usage_count")
	} else {
		db = db.Select("t.external_id as template_external_id, t.name, t.tags, t.usage_count as usage_count, t.usage_count as global_usage_count")
	}

	if err := db.Order("t.name asc").Order("t.id asc").Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]TemplateLibraryItem, 0, len(rows))
	for _, it := range rows {
		out = append(out, TemplateLibraryItem{
			TemplateExternalID: it.TemplateExternalID,
			Name:               it.Name,
			Tags:               it.Tags,
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
