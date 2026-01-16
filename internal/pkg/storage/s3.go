package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader struct {
	cli      *s3.Client
	bucket   string
	endpoint string
}

func NewS3(bucket, region, endpoint, accessKey, secretKey string) (Uploader, error) {
	opts := []func(*awscfg.LoadOptions) error{awscfg.WithRegion(region)}
	if accessKey != "" && secretKey != "" {
		opts = append(opts, awscfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")))
	}
	cfg, err := awscfg.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	cli := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
			// o.UsePathStyle = true
		}
	})
	return &S3Uploader{cli: cli, bucket: bucket, endpoint: endpoint}, nil
}

func (s *S3Uploader) Upload(ctx context.Context, name string, content []byte) (string, error) {
	_, err := s.cli.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(name),
		Body:        bytes.NewReader(content),
		ContentType: aws.String("application/octet-stream"),
	})
	if err != nil {
		return "", err
	}

	if s.endpoint != "" {
		ep := s.endpoint
		scheme := "https"
		if strings.HasPrefix(ep, "http://") {
			scheme = "http"
			ep = strings.TrimPrefix(ep, "http://")
		} else if strings.HasPrefix(ep, "https://") {
			scheme = "https"
			ep = strings.TrimPrefix(ep, "https://")
		}
		ep = strings.TrimRight(ep, "/")
		return scheme + "://" + s.bucket + "." + ep + "/" + name, nil
	}

	return "https://" + s.bucket + ".s3.amazonaws.com/" + name, nil
}

func (s *S3Uploader) Delete(ctx context.Context, urlStr string) error {
	if urlStr == "" {
		return nil
	}
	key := path.Base(urlStr)
	if key == "" || key == "/" || key == "." {
		return nil
	}
	_, _ = s.cli.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return nil
}

func (s *S3Uploader) Download(ctx context.Context, urlStr string) (io.ReadCloser, error) {
	key := path.Base(urlStr)
	if key == "" || key == "/" || key == "." {
		return nil, fmt.Errorf("invalid key")
	}
	out, err := s.cli.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}
