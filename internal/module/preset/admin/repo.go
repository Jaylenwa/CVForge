package admin

import presetmod "cvforge/internal/module/preset"

type Repo struct {
	inner *presetmod.Repo
}

func DefaultRepo() *Repo {
	return &Repo{inner: presetmod.DefaultRepo()}
}

func (r *Repo) AdminListContentPresets(page, size int, q, role, language string) ([]presetmod.ContentPreset, int64, error) {
	return r.inner.AdminListContentPresets(page, size, q, role, language)
}

func (r *Repo) AdminCreateContentPreset(p *presetmod.ContentPreset) error {
	return r.inner.AdminCreateContentPreset(p)
}

func (r *Repo) AdminPatchContentPreset(id uint, patch map[string]any) error {
	return r.inner.AdminPatchContentPreset(id, patch)
}

func (r *Repo) AdminDeleteContentPreset(id uint) error {
	return r.inner.AdminDeleteContentPreset(id)
}

