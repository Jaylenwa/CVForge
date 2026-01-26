package library

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

func (s *Service) ListTemplateVariants(roleID uint, categoryID uint, q string) ([]TemplateVariant, error) {
	return s.repo.ListTemplateVariants(roleID, categoryID, q)
}
