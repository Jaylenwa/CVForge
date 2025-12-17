package upload

import (
	"openresume/internal/infra/config"
	"openresume/internal/pkg/storage"

	"github.com/gin-gonic/gin"
)

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Upload(c *gin.Context, name string, b []byte) (string, error) {
	up, e := storage.New(config.Load())
	if e != nil {
		return "", e
	}
	return up.Upload(c, name, b)
}
