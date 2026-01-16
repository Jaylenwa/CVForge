package upload

import (
	conf "openresume/internal/module/config"
	"openresume/internal/pkg/storage"

	"github.com/gin-gonic/gin"
)

type Service struct {
	sys *conf.Service
}

func NewService() *Service {
	return &Service{
		sys: conf.NewService(),
	}
}

func (s *Service) Upload(c *gin.Context, name string, b []byte) (string, error) {
	cfg := s.sys.GetStorageSettings()
	up, e := storage.NewFromSettings(cfg)
	if e != nil {
		return "", e
	}
	return up.Upload(c, name, b)
}
