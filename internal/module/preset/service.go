package preset

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

func (s *Service) GetByIDActive(id uint) (ContentPreset, error) {
	return s.repo.GetByIDActive(id)
}
