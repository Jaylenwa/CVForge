package admin

import taxmod "cvforge/internal/module/taxonomy"

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

func (s *Service) AdminListJobCategories(page, size int, q string, parentID *uint) ([]taxmod.JobCategory, int64, error) {
	return s.repo.AdminListJobCategories(page, size, q, parentID)
}

func (s *Service) ListJobCategoryI18n(categoryIDs []uint) ([]taxmod.JobCategoryI18n, error) {
	return s.repo.ListJobCategoryI18n(categoryIDs)
}

func (s *Service) AdminCreateJobCategoryWithNames(m *taxmod.JobCategory, names map[string]string) error {
	return s.repo.AdminCreateJobCategoryWithNames(m, names)
}

func (s *Service) AdminPatchJobCategoryWithNames(id uint, patch map[string]any, names map[string]string) error {
	return s.repo.AdminPatchJobCategoryWithNames(id, patch, names)
}

func (s *Service) AdminDeleteJobCategory(id uint) error {
	return s.repo.AdminDeleteJobCategory(id)
}

func (s *Service) AdminListJobRoles(page, size int, q string, categoryID *uint) ([]taxmod.JobRole, int64, error) {
	return s.repo.AdminListJobRoles(page, size, q, categoryID)
}

func (s *Service) ListJobRoleI18n(roleIDs []uint) ([]taxmod.JobRoleI18n, error) {
	return s.repo.ListJobRoleI18n(roleIDs)
}

func (s *Service) AdminCreateJobRoleWithNames(m *taxmod.JobRole, names map[string]string) error {
	return s.repo.AdminCreateJobRoleWithNames(m, names)
}

func (s *Service) AdminPatchJobRoleWithNames(id uint, patch map[string]any, names map[string]string) error {
	return s.repo.AdminPatchJobRoleWithNames(id, patch, names)
}

func (s *Service) AdminDeleteJobRole(id uint) error {
	return s.repo.AdminDeleteJobRole(id)
}

