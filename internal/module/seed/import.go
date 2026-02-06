package seed

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"openresume/internal/module/preset"
	"openresume/internal/module/taxonomy"

	"gorm.io/gorm"
)

func Import(ctx context.Context, db *gorm.DB, s SeedData) (ImportCounts, error) {
	if db == nil {
		return ImportCounts{}, errors.New("db is nil")
	}
	if err := s.Validate(); err != nil {
		return ImportCounts{}, err
	}

	var counts ImportCounts
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		taxRepo := taxonomy.NewRepo(tx)
		presetRepo := preset.NewRepo(tx)

		categoryIDByExternal, err := importCategories(tx, taxRepo, s.Categories, &counts)
		if err != nil {
			return err
		}

		roleIDByExternal := make(map[string]uint, len(s.Roles))
		for _, sr := range s.Roles {
			cid := categoryIDByExternal[sr.CategoryExternalID]
			externalID := strings.TrimSpace(sr.ExternalID)
			m := taxonomy.JobRole{
				ExternalID: &externalID,
				CategoryID: cid,
				OrderNum:   sr.OrderNum,
				IsActive:   sr.IsActive,
			}
			if err := taxRepo.UpsertJobRoleWithNames(tx, &m, sr.Names); err != nil {
				return fmt.Errorf("upsert role %q: %w", sr.ExternalID, err)
			}
			roleIDByExternal[sr.ExternalID] = m.ID
			counts.JobRoles++
		}

		for _, sp := range s.Presets {
			rid, ok := roleIDByExternal[sp.RoleCode]
			if !ok || rid == 0 {
				return fmt.Errorf("preset %q references missing role %q", sp.ExternalID, sp.RoleCode)
			}
			m := preset.ContentPreset{
				ExternalID: sp.ExternalID,
				Name:       sp.Name,
				Language:   sp.Language,
				RoleID:     rid,
				DataJSON:   sp.DataJSON,
				IsActive:   sp.IsActive,
			}
			if err := presetRepo.UpsertContentPreset(tx, &m); err != nil {
				return fmt.Errorf("upsert preset %q: %w", sp.ExternalID, err)
			}
			counts.ContentPresets++
		}
		return nil
	})
	if err != nil {
		return ImportCounts{}, err
	}
	return counts, nil
}

func importCategories(tx *gorm.DB, repo *taxonomy.Repo, list []SeedJobCategory, counts *ImportCounts) (map[string]uint, error) {
	categoryIDByExternal := make(map[string]uint, len(list))
	pending := make(map[string]SeedJobCategory, len(list))
	for _, c := range list {
		pending[c.ExternalID] = c
	}

	for len(pending) > 0 {
		progress := 0
		for externalID, sc := range pending {
			if strings.TrimSpace(sc.ParentExternalID) != "" {
				if _, ok := categoryIDByExternal[sc.ParentExternalID]; !ok {
					continue
				}
			}
			var parentID *uint
			if strings.TrimSpace(sc.ParentExternalID) != "" {
				pid := categoryIDByExternal[sc.ParentExternalID]
				parentID = &pid
			}
			id := strings.TrimSpace(sc.ExternalID)
			m := taxonomy.JobCategory{
				ExternalID: &id,
				ParentID:   parentID,
				OrderNum:   sc.OrderNum,
				IsActive:   sc.IsActive,
			}
			if err := repo.UpsertJobCategoryWithNames(tx, &m, sc.Names); err != nil {
				return nil, fmt.Errorf("upsert category %q: %w", sc.ExternalID, err)
			}
			categoryIDByExternal[sc.ExternalID] = m.ID
			(*counts).JobCategories++
			delete(pending, externalID)
			progress++
		}
		if progress == 0 {
			var waiting []string
			for k := range pending {
				waiting = append(waiting, k)
			}
			return nil, fmt.Errorf("cannot resolve category parents for: %s", strings.Join(waiting, ", "))
		}
	}
	return categoryIDByExternal, nil
}

