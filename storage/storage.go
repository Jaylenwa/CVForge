package storage

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"openresume/config"
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
