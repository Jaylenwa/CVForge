package upload

import (
	"context"

	conf "cvforge/internal/module/config"
	"cvforge/internal/pkg/storage"
)

type Service struct {
	sys *conf.Service
}

func NewService() *Service {
	return &Service{
		sys: conf.NewService(),
	}
}

func (s *Service) Upload(ctx context.Context, name string, b []byte) (string, error) {
	cfg := s.sys.GetStorageSettings()
	up, e := storage.NewFromSettings(cfg)
	if e != nil {
		return "", e
	}
	return up.Upload(ctx, name, b)
}
