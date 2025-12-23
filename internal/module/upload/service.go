package upload

import (
	"openresume/internal/infra/config"
	conf "openresume/internal/module/config"
	"openresume/internal/pkg/storage"

	"github.com/gin-gonic/gin"
)

type Service struct {
	cfg config.Config
	sys *conf.Service
}

func NewService(cfg config.Config, sys *conf.Service) *Service { return &Service{cfg: cfg, sys: sys} }

func (s *Service) Upload(c *gin.Context, name string, b []byte) (string, error) {
	useS3 := s.sys.GetBool("enabled_storage_s3", false)
	bucket := s.sys.GetWithDefault("storage_s3_bucket", "")
	region := s.sys.GetWithDefault("storage_s3_region", "")
	endpoint := s.sys.GetWithDefault("storage_s3_endpoint", "")
	accessKey := s.sys.GetWithDefault("storage_s3_access_key", "")
	secretKey := s.sys.GetWithDefault("storage_s3_secret_key", "")
	up, e := storage.New(useS3, bucket, region, endpoint, accessKey, secretKey)
	if e != nil {
		return "", e
	}
	return up.Upload(c, name, b)
}
