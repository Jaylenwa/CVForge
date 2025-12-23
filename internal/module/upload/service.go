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
	useS3 := s.sys.GetBool("storage_s3_enabled", false)
	bucket := s.sys.GetWithDefault("storage_s3_bucket", s.cfg.S3Bucket)
	region := s.sys.GetWithDefault("storage_s3_region", s.cfg.S3Region)
	endpoint := s.sys.GetWithDefault("storage_s3_endpoint", s.cfg.S3Endpoint)
	accessKey := s.sys.GetWithDefault("storage_s3_access_key", s.cfg.S3AccessKey)
	secretKey := s.sys.GetWithDefault("storage_s3_secret_key", s.cfg.S3SecretKey)
	up, e := storage.NewFromValues(useS3, bucket, region, endpoint, accessKey, secretKey)
	if e != nil {
		return "", e
	}
	return up.Upload(c, name, b)
}
