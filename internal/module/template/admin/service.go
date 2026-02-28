package admin

import (
	"context"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"
	templatemod "cvforge/internal/module/template"
)

type Service struct {
	repo *templatemod.Repo
}

func NewService() *Service {
	return &Service{repo: templatemod.DefaultRepo()}
}

func (s *Service) Create(externalID string, names map[string]string) error {
	t := templatemod.Template{ExternalID: externalID}
	if err := s.repo.CreateWithNames(&t, names); err != nil {
		return err
	}
	cache.RDB.Del(context.Background(), string(common.RedisKeyTemplatesListAll))
	return nil
}

func (s *Service) UpdateNames(id string, names map[string]string) error {
	t, err := s.repo.GetByExternal(id)
	if err != nil {
		return err
	}
	if err := s.repo.PatchWithNames(t.ID, nil, names); err != nil {
		return err
	}
	cache.RDB.Del(context.Background(), string(common.RedisKeyTemplatesListAll))
	return nil
}

func (s *Service) Delete(id string) error {
	if err := s.repo.DeleteByExternal(id); err != nil {
		return err
	}
	cache.RDB.Del(context.Background(), string(common.RedisKeyTemplatesListAll))
	return nil
}
