package catalog

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
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

