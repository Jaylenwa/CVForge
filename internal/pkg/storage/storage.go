package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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

func New(useS3 bool, bucket, region, endpoint, accessKey, secretKey string) (Uploader, error) {
	if useS3 {
		if bucket == "" || region == "" {
			return nil, fmt.Errorf("s3 bucket and region required")
		}
		return NewS3(bucket, region, endpoint, accessKey, secretKey)
	}
	return Local{}, nil
}
