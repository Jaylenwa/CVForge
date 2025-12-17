package storage

import (
	"bytes"
	"context"
	"openresume/internal/infra/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader struct {
	cli    *s3.Client
	bucket string
}

func NewS3(cfgApp config.Config) (Uploader, error) {
	opts := []func(*awscfg.LoadOptions) error{awscfg.WithRegion(cfgApp.S3Region)}
	if cfgApp.S3AccessKey != "" && cfgApp.S3SecretKey != "" {
		opts = append(opts, awscfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfgApp.S3AccessKey, cfgApp.S3SecretKey, "")))
	}
	if cfgApp.S3Endpoint != "" {
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: cfgApp.S3Endpoint, HostnameImmutable: true}, nil
		})
		opts = append(opts, awscfg.WithEndpointResolverWithOptions(resolver))
	}
	cfg, err := awscfg.LoadDefaultConfig(context.Background(), opts...)
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
