package catalog

import "gorm.io/gorm"

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

type ImportCounts struct {
	JobCategories int `json:"jobCategories"`
	JobRoles      int `json:"jobRoles"`
	ContentPresets int `json:"contentPresets"`
	TemplateVariants int `json:"templateVariants"`
}

func (s *Service) ListJobCategories() ([]JobCategory, error) {
	return s.repo.ListJobCategories()
}

func (s *Service) ListJobRoles(categoryExternalID string, q string) ([]JobRole, error) {
	return s.repo.ListJobRoles(categoryExternalID, q)
}

func (s *Service) ListTemplateVariants(roleExternalID string, categoryExternalID string, q string) ([]TemplateVariant, error) {
	return s.repo.ListTemplateVariants(roleExternalID, categoryExternalID, q)
}

func (s *Service) GetContentPresetByExternal(id string) (ContentPreset, error) {
	return s.repo.GetContentPresetByExternal(id)
}

func (s *Service) ImportSeed(seed SeedData) (ImportCounts, error) {
	var counts ImportCounts
	err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		for i := range seed.Categories {
			if err := s.repo.UpsertJobCategory(tx, &seed.Categories[i]); err != nil {
				return err
			}
			counts.JobCategories++
		}
		for i := range seed.Roles {
			if err := s.repo.UpsertJobRole(tx, &seed.Roles[i]); err != nil {
				return err
			}
			counts.JobRoles++
		}
		for i := range seed.Presets {
			if err := s.repo.UpsertContentPreset(tx, &seed.Presets[i]); err != nil {
				return err
			}
			counts.ContentPresets++
		}
		for i := range seed.Variants {
			if err := s.repo.UpsertTemplateVariant(tx, &seed.Variants[i]); err != nil {
				return err
			}
			counts.TemplateVariants++
		}
		return nil
	})
	return counts, err
}
