package admin

import taxmod "cvforge/internal/module/taxonomy"

type Repo struct {
	inner *taxmod.Repo
}

func DefaultRepo() *Repo {
	return &Repo{inner: taxmod.DefaultRepo()}
}

func (r *Repo) AdminListJobCategories(page, size int, q string, parentID *uint) ([]taxmod.JobCategory, int64, error) {
	return r.inner.AdminListJobCategories(page, size, q, parentID)
}

func (r *Repo) ListJobCategoryI18n(categoryIDs []uint) ([]taxmod.JobCategoryI18n, error) {
	return r.inner.ListJobCategoryI18n(categoryIDs)
}

func (r *Repo) AdminCreateJobCategoryWithNames(m *taxmod.JobCategory, names map[string]string) error {
	return r.inner.AdminCreateJobCategoryWithNames(m, names)
}

func (r *Repo) AdminPatchJobCategoryWithNames(id uint, patch map[string]any, names map[string]string) error {
	return r.inner.AdminPatchJobCategoryWithNames(id, patch, names)
}

func (r *Repo) AdminDeleteJobCategory(id uint) error {
	return r.inner.AdminDeleteJobCategory(id)
}

func (r *Repo) AdminListJobRoles(page, size int, q string, categoryID *uint) ([]taxmod.JobRole, int64, error) {
	return r.inner.AdminListJobRoles(page, size, q, categoryID)
}

func (r *Repo) ListJobRoleI18n(roleIDs []uint) ([]taxmod.JobRoleI18n, error) {
	return r.inner.ListJobRoleI18n(roleIDs)
}

func (r *Repo) AdminCreateJobRoleWithNames(m *taxmod.JobRole, names map[string]string) error {
	return r.inner.AdminCreateJobRoleWithNames(m, names)
}

func (r *Repo) AdminPatchJobRoleWithNames(id uint, patch map[string]any, names map[string]string) error {
	return r.inner.AdminPatchJobRoleWithNames(id, patch, names)
}

func (r *Repo) AdminDeleteJobRole(id uint) error {
	return r.inner.AdminDeleteJobRole(id)
}

