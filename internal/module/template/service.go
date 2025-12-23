package template

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repo
	rdb  *redis.Client
}

func NewService(db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{repo: NewRepo(db), rdb: rdb}
}

func (s *Service) ListAllPayload() (string, error) {
	if s.rdb != nil {
		if val, err := s.rdb.Get(context.Background(), "templates:list:all").Result(); err == nil {
			return val, nil
		}
	}
	list, err := s.repo.ListAll()
	if err != nil {
		return "", err
	}
	payloadBytes, _ := json.Marshal(map[string]any{"items": list})
	payload := string(payloadBytes)
	if s.rdb != nil {
		_ = s.rdb.Set(context.Background(), "templates:list:all", payload, time.Hour).Err()
	}
	return payload, nil
}

func (s *Service) GetByExternal(id string) (Template, error) {
	return s.repo.GetByExternal(id)
}

func (s *Service) Create(t Template) error {
	if err := s.repo.Create(&t); err != nil {
		return err
	}
	if s.rdb != nil {
		_ = s.rdb.Del(context.Background(), "templates:list:all").Err()
	}
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
	if s.rdb != nil {
		_ = s.rdb.Del(context.Background(), "templates:list:all").Err()
	}
	return nil
}

func (s *Service) Delete(id string) error {
	if err := s.repo.DeleteByExternal(id); err != nil {
		return err
	}
	if s.rdb != nil {
		_ = s.rdb.Del(context.Background(), "templates:list:all").Err()
	}
	return nil
}
