package seed

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"cvforge/internal/infra/database"
	"cvforge/internal/module/preset"
	"cvforge/internal/module/taxonomy"
	tplmod "cvforge/internal/module/template"
	"cvforge/internal/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SeedData struct {
	Categories []SeedJobCategory   `json:"Categories"`
	Roles      []SeedJobRole       `json:"Roles"`
	Presets    []SeedContentPreset `json:"Presets"`
}

type SeedJobCategory struct {
	ExternalID       string            `json:"ExternalID"`
	Names            map[string]string `json:"Names"`
	ParentExternalID string            `json:"ParentExternalID"`
	OrderNum         int               `json:"OrderNum"`
	IsActive         bool              `json:"IsActive"`
}

type SeedJobRole struct {
	ExternalID         string            `json:"ExternalID"`
	CategoryExternalID string            `json:"CategoryExternalID"`
	Names              map[string]string `json:"Names"`
	OrderNum           int               `json:"OrderNum"`
	IsActive           bool              `json:"IsActive"`
}

type SeedContentPreset struct {
	ExternalID string `json:"ExternalID"`
	Name       string `json:"Name"`
	Language   string `json:"Language"`
	RoleCode   string `json:"RoleCode"`
	DataJSON   string `json:"DataJSON"`
	IsActive   bool   `json:"IsActive"`
}

type ImportCounts struct {
	JobCategories  int `json:"jobCategories"`
	JobRoles       int `json:"jobRoles"`
	ContentPresets int `json:"contentPresets"`
	Templates      int `json:"templates"`
}

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

func LoadFromBytes(b []byte) (SeedData, error) {
	var s SeedData
	if err := json.Unmarshal(b, &s); err != nil {
		return SeedData{}, fmt.Errorf("parse seed json: %w", err)
	}
	if err := s.Validate(); err != nil {
		return SeedData{}, err
	}
	return s, nil
}

func loadFromBytesNoValidate(b []byte) (SeedData, error) {
	var s SeedData
	if err := json.Unmarshal(b, &s); err != nil {
		return SeedData{}, fmt.Errorf("parse seed json: %w", err)
	}
	return s, nil
}

func LoadFromFS(fsys fs.FS, name string) (SeedData, error) {
	b, err := fs.ReadFile(fsys, name)
	if err != nil {
		return SeedData{}, fmt.Errorf("read seed file: %w", err)
	}
	return LoadFromBytes(b)
}

func loadFromFSNoValidate(fsys fs.FS, name string) (SeedData, error) {
	b, err := fs.ReadFile(fsys, name)
	if err != nil {
		return SeedData{}, fmt.Errorf("read seed file: %w", err)
	}
	return loadFromBytesNoValidate(b)
}

func LoadFromFilePath(path string) (SeedData, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return SeedData{}, fmt.Errorf("read seed file: %w", err)
	}
	return LoadFromBytes(b)
}

func LoadFromDirPath(dir string) (SeedData, error) {
	parts := []string{
		"categories.json",
		"roles.json",
		"presets.json",
	}

	var merged SeedData
	for _, name := range parts {
		path := filepath.Join(dir, name)
		b, err := os.ReadFile(path)
		if err != nil {
			return SeedData{}, fmt.Errorf("read seed file: %w", err)
		}
		part, err := loadFromBytesNoValidate(b)
		if err != nil {
			return SeedData{}, err
		}
		merged.Categories = append(merged.Categories, part.Categories...)
		merged.Roles = append(merged.Roles, part.Roles...)
		merged.Presets = append(merged.Presets, part.Presets...)
	}

	if err := merged.Validate(); err != nil {
		return SeedData{}, err
	}
	return merged, nil
}

//go:embed default/categories.json default/roles.json default/presets.json
var defaultFS embed.FS

//go:embed default/templates.json
var defaultTemplatesFS embed.FS

func LoadDefaultSeed() (SeedData, error) {
	parts := []string{
		"default/categories.json",
		"default/roles.json",
		"default/presets.json",
	}

	var merged SeedData
	for _, name := range parts {
		part, err := loadFromFSNoValidate(defaultFS, name)
		if err != nil {
			return SeedData{}, err
		}
		merged.Categories = append(merged.Categories, part.Categories...)
		merged.Roles = append(merged.Roles, part.Roles...)
		merged.Presets = append(merged.Presets, part.Presets...)
	}

	if err := merged.Validate(); err != nil {
		return SeedData{}, err
	}
	return merged, nil
}

func LoadDefaultTemplateItems() ([]tplmod.SeedTemplateItem, error) {
	b, err := defaultTemplatesFS.ReadFile("default/templates.json")
	if err != nil {
		return nil, fmt.Errorf("read default templates: %w", err)
	}
	var wrapper struct {
		Templates []tplmod.SeedTemplateItem `json:"Templates"`
	}
	if err := json.Unmarshal(b, &wrapper); err == nil && len(wrapper.Templates) > 0 {
		return wrapper.Templates, nil
	}
	var arr []tplmod.SeedTemplateItem
	if err := json.Unmarshal(b, &arr); err == nil {
		return arr, nil
	}
	return nil, fmt.Errorf("invalid default templates json")
}
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

func RunCLI(args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "usage: cvforge seed <import-default|import|templates-import> [--file path]")
		return 2
	}

	switch args[0] {
	case "import-default":
		s, err := LoadDefaultSeed()
		if err != nil {
			logger.L().Error("load default seed failed", zap.Error(err))
			return 1
		}
		counts, err := Import(context.Background(), database.DB, s)
		if err != nil {
			logger.L().Error("seed import failed", zap.Error(err))
			return 1
		}
		// Also import default templates to align with frontend template IDs
		if items, terr := LoadDefaultTemplateItems(); terr == nil && len(items) > 0 {
			if ierr := tplmod.NewService().Seed(items); ierr != nil {
				logger.L().Error("templates default import failed", zap.Error(ierr))
				return 1
			}
			counts.Templates = len(items)
		} else if terr != nil {
			logger.L().Error("load default templates failed", zap.Error(terr))
			return 1
		}
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"success": true, "counts": counts})
		return 0

	case "import":
		fs := flag.NewFlagSet("cvforge seed import", flag.ContinueOnError)
		fs.SetOutput(os.Stderr)
		filePath := fs.String("file", "", "seed json file path")
		dirPath := fs.String("dir", "", "seed directory path (categories.json, roles.json, presets.json)")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		if (*filePath == "" && *dirPath == "") || (*filePath != "" && *dirPath != "") {
			_, _ = fmt.Fprintln(os.Stderr, "usage: cvforge seed import --file path/to/seed.json | --dir path/to/seed-dir")
			return 2
		}
		var s SeedData
		var err error
		if *dirPath != "" {
			s, err = LoadFromDirPath(*dirPath)
		} else {
			s, err = LoadFromFilePath(*filePath)
		}
		if err != nil {
			logger.L().Error("load seed file failed", zap.Error(err))
			return 1
		}
		counts, err := Import(context.Background(), database.DB, s)
		if err != nil {
			logger.L().Error("seed import failed", zap.Error(err))
			return 1
		}
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"success": true, "counts": counts})
		return 0

	default:
		_, _ = fmt.Fprintln(os.Stderr, "usage: cvforge seed <import-default|import|templates-import|templates-import-default> [--file path]")
		return 2
	}
}
