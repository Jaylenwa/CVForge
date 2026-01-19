package preset

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

func (s *Service) GetByExternalActive(id string) (ContentPreset, error) {
	return s.repo.GetByExternalActive(id)
}

