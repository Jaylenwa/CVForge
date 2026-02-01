package library

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

func (s *Service) ListTemplateLibraryItems(roleID uint, language, sort string) ([]TemplateLibraryItem, error) {
	return s.repo.ListTemplateLibraryItems(roleID, language, sort)
}
