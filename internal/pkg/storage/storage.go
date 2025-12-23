package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"openresume/internal/infra/config"
)

type Uploader interface {
	Upload(ctx context.Context, name string, content []byte) (string, error)
}

type Local struct{}

func (Local) Upload(ctx context.Context, name string, content []byte) (string, error) {
	_ = os.MkdirAll("uploads", 0755)
	p := filepath.Join("uploads", name)
	if err := os.WriteFile(p, content, 0644); err != nil {
		return "", err
	}
	return "/public/uploads/" + name, nil
}

func New(cfg config.Config) (Uploader, error) {
	if strings.ToLower(cfg.UploadBackend) == "s3" {
		if cfg.S3Bucket == "" {
			return Local{}, nil
		}
		return NewS3(cfg)
	}
	return Local{}, nil
}

func NewFromValues(useS3 bool, bucket, region, endpoint, accessKey, secretKey string) (Uploader, error) {
	if useS3 {
		if bucket == "" || region == "" {
			return nil, fmt.Errorf("s3 bucket and region required")
		}
		c := config.Config{
			S3Bucket:    bucket,
			S3Region:    region,
			S3Endpoint:  endpoint,
			S3AccessKey: accessKey,
			S3SecretKey: secretKey,
		}
		return NewS3(c)
	}
	return Local{}, nil
}
