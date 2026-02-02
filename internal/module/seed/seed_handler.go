package seed

import (
	"net/http"

	"openresume/internal/infra/database"
	"openresume/internal/module/preset"
	"openresume/internal/module/seed/presets"
	"openresume/internal/module/taxonomy"
	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SeedData struct {
	Categories []SeedJobCategory
	Roles      []SeedJobRole
	Presets    []SeedContentPreset
}

type SeedJobCategory struct {
	ExternalID       string
	Name             string
	ParentExternalID string
	OrderNum         int
	IsActive         bool
}

type SeedJobRole struct {
	ExternalID         string
	CategoryExternalID string
	Name               string
	Tags               string
	OrderNum           int
	IsActive           bool
}

type SeedContentPreset struct {
	Name     string
	Language string
	RoleCode string
	Tags     string
	DataJSON string
	IsActive bool
}

type ImportCounts struct {
	JobCategories  int `json:"jobCategories"`
	JobRoles       int `json:"jobRoles"`
	ContentPresets int `json:"contentPresets"`
}

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) AdminImportDefault(c *gin.Context) {
	seed, err := DefaultSeed()
	if err != nil {
		logger.WithCtx(c).Error("seed.build failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "seed error"})
		return
	}

	var counts ImportCounts
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		taxRepo := taxonomy.NewRepo(tx)
		presetRepo := preset.NewRepo(tx)

		categoryIDByExternal := make(map[string]uint, len(seed.Categories))
		for _, sc := range seed.Categories {
			var parentID *uint
			if sc.ParentExternalID != "" {
				pid, ok := categoryIDByExternal[sc.ParentExternalID]
				if !ok || pid == 0 {
					return gorm.ErrInvalidData
				}
				parentID = &pid
			}
			m := taxonomy.JobCategory{
				Name:     sc.Name,
				ParentID: parentID,
				OrderNum: sc.OrderNum,
				IsActive: sc.IsActive,
			}
			if err := taxRepo.UpsertJobCategory(tx, &m); err != nil {
				return err
			}
			categoryIDByExternal[sc.ExternalID] = m.ID
			counts.JobCategories++
		}

		roleIDByExternal := make(map[string]uint, len(seed.Roles))
		for _, sr := range seed.Roles {
			cid, ok := categoryIDByExternal[sr.CategoryExternalID]
			if !ok || cid == 0 {
				return gorm.ErrInvalidData
			}
			m := taxonomy.JobRole{
				CategoryID: cid,
				Name:       sr.Name,
				Tags:       sr.Tags,
				OrderNum:   sr.OrderNum,
				IsActive:   sr.IsActive,
			}
			if err := taxRepo.UpsertJobRole(tx, &m); err != nil {
				return err
			}
			roleIDByExternal[sr.ExternalID] = m.ID
			counts.JobRoles++
		}

		for _, sp := range seed.Presets {
			rid, ok := roleIDByExternal[sp.RoleCode]
			if !ok || rid == 0 {
				return gorm.ErrInvalidData
			}
			m := preset.ContentPreset{
				Name:     sp.Name,
				Language: sp.Language,
				RoleID:   rid,
				Tags:     sp.Tags,
				DataJSON: sp.DataJSON,
				IsActive: sp.IsActive,
			}
			if err := presetRepo.UpsertContentPreset(tx, &m); err != nil {
				return err
			}
			counts.ContentPresets++
		}
		return nil
	})
	if err != nil {
		logger.WithCtx(c).Error("seed.import failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "counts": counts})
}

func DefaultSeed() (SeedData, error) {
	return SeedData{
		Categories: SeedCategories,
		Roles:      SeedRoles,
		Presets: []SeedContentPreset{
			{Name: "Java 开发（中文示例）", Language: "zh", RoleCode: "Java", Tags: "Java,后端,中文", DataJSON: string(presets.GenerateJavaPreset()), IsActive: true},
			{Name: "Python 开发（中文示例）", Language: "zh", RoleCode: "Python", Tags: "Python,后端,中文", DataJSON: string(presets.GeneratePythonPreset()), IsActive: true},
			{Name: "Go 开发（中文示例）", Language: "zh", RoleCode: "golang", Tags: "Go,后端,中文", DataJSON: string(presets.GenerateGolangPreset()), IsActive: true},
		},
	}, nil
}
