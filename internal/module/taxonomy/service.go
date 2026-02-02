package taxonomy

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

func (s *Service) ListJobCategories(language string) ([]JobCategoryView, error) {
	return s.repo.ListJobCategories(language)
}

func (s *Service) ListJobRoles(language string, categoryID uint, q string) ([]JobRoleView, error) {
	return s.repo.ListJobRoles(language, categoryID, q)
}
