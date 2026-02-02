package preset

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

func (s *Service) GetByIDActive(id uint, language string) (ContentPreset, error) {
	return s.repo.GetByIDActive(id, language)
}
