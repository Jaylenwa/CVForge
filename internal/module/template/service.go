package template

import (
	"context"
	"encoding/json"
	"time"

	"openresume/internal/common"
	"openresume/internal/infra/cache"
)

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

func (s *Service) ListAllPayload() (string, error) {
	if val, err := cache.RDB.Get(context.Background(), string(common.RedisKeyTemplatesListAll)).Result(); err == nil {
		return val, nil
	}
	list, err := s.repo.ListAll()
	if err != nil {
		return "", err
	}
	payloadBytes, _ := json.Marshal(map[string]any{"items": list})
	payload := string(payloadBytes)
	cache.RDB.Set(context.Background(), string(common.RedisKeyTemplatesListAll), payload, time.Hour)
	return payload, nil
}

func (s *Service) GetByExternal(id string) (Template, error) {
	return s.repo.GetByExternal(id)
}

func (s *Service) Create(t Template) error {
	if err := s.repo.Create(&t); err != nil {
		return err
	}
	cache.RDB.Del(context.Background(), string(common.RedisKeyTemplatesListAll))
	return nil
}

func (s *Service) Update(id string, patch func(*Template)) error {
	t, err := s.repo.GetByExternal(id)
	if err != nil {
		return err
	}
	patch(&t)
	if err := s.repo.Save(&t); err != nil {
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
