package taxonomy

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

