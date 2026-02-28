package admin

import presetmod "cvforge/internal/module/preset"

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

func (s *Service) AdminListContentPresets(page, size int, q, role, language string) ([]presetmod.ContentPreset, int64, error) {
	return s.repo.AdminListContentPresets(page, size, q, role, language)
}

func (s *Service) AdminCreateContentPreset(p *presetmod.ContentPreset) error {
	return s.repo.AdminCreateContentPreset(p)
}

func (s *Service) AdminPatchContentPreset(id uint, patch map[string]any) error {
	return s.repo.AdminPatchContentPreset(id, patch)
}

func (s *Service) AdminDeleteContentPreset(id uint) error {
	return s.repo.AdminDeleteContentPreset(id)
}

