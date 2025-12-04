package storage

import (
	"bytes"
	"context"
	"openresume/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader struct {
	cli    *s3.Client
	bucket string
}

func NewS3(cfgApp config.Config) (Uploader, error) {
	cfg, err := awscfg.LoadDefaultConfig(context.Background(), awscfg.WithRegion(cfgApp.S3Region))
	if err != nil {
		return nil, err
	}
	cli := s3.NewFromConfig(cfg)
	return &S3Uploader{cli: cli, bucket: cfgApp.S3Bucket}, nil
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
	return "https://" + s.bucket + ".s3.amazonaws.com/" + name, nil
}
