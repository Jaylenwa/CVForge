package upload

import (
	"openresume/internal/common"
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
	useS3 := s.sys.GetBool(string(common.ConfigKeyEnabledStorageS3), false)
	bucket := s.sys.Get(string(common.ConfigKeyStorageS3Bucket))
	region := s.sys.Get(string(common.ConfigKeyStorageS3Region))
	endpoint := s.sys.Get(string(common.ConfigKeyStorageS3Endpoint))
	accessKey := s.sys.Get(string(common.ConfigKeyStorageS3AccessKey))
	secretKey := s.sys.Get(string(common.ConfigKeyStorageS3SecretKey))
	up, e := storage.New(useS3, bucket, region, endpoint, accessKey, secretKey)
	if e != nil {
		return "", e
	}
	return up.Upload(c, name, b)
}
