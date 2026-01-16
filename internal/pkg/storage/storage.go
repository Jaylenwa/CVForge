package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Uploader interface {
	Upload(ctx context.Context, name string, content []byte) (string, error)
	Delete(ctx context.Context, url string) error
	Download(ctx context.Context, url string) (io.ReadCloser, error)
}

type Settings struct {
	Enabled   bool
	Bucket    string
	Region    string
	Endpoint  string
	AccessKey string
	SecretKey string
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

func (Local) Delete(ctx context.Context, urlStr string) error {
	if urlStr == "" {
		return nil
	}
	u, err := url.Parse(urlStr)
	var pth string
	if err == nil && u.Scheme != "" {
		pth = u.Path
	} else {
		pth = urlStr
	}
	var name string
	if strings.HasPrefix(pth, "/public/uploads/") {
		name = strings.TrimPrefix(pth, "/public/uploads/")
	} else if strings.HasPrefix(pth, "uploads/") {
		name = strings.TrimPrefix(pth, "uploads/")
	} else {
		name = path.Base(pth)
	}
	if name == "" || name == "/" || name == "." {
		return nil
	}
	fp := filepath.Join("uploads", name)
	_ = os.Remove(fp)
	return nil
}

func (Local) Download(ctx context.Context, urlStr string) (io.ReadCloser, error) {
	if urlStr == "" {
		return nil, fmt.Errorf("empty url")
	}
	u, err := url.Parse(urlStr)
	var pth string
	if err == nil && u.Scheme != "" {
		pth = u.Path
	} else {
		pth = urlStr
	}
	var name string
	if strings.HasPrefix(pth, "/public/uploads/") {
		name = strings.TrimPrefix(pth, "/public/uploads/")
	} else if strings.HasPrefix(pth, "uploads/") {
		name = strings.TrimPrefix(pth, "uploads/")
	} else {
		name = path.Base(pth)
	}
	if name == "" || name == "/" || name == "." {
		return nil, fmt.Errorf("invalid name")
	}
	fp := filepath.Join("uploads", name)
	return os.Open(fp)
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

func NewFromSettings(s Settings) (Uploader, error) {
	return New(s.Enabled, s.Bucket, s.Region, s.Endpoint, s.AccessKey, s.SecretKey)
}
