package storage

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader struct {
	cli    *s3.Client
	bucket string
}

func NewS3(bucket, region, endpoint, accessKey, secretKey string) (Uploader, error) {
	opts := []func(*awscfg.LoadOptions) error{awscfg.WithRegion(region)}
	if accessKey != "" && secretKey != "" {
		opts = append(opts, awscfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")))
	}
	if endpoint != "" {
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: endpoint, HostnameImmutable: true}, nil
		})
		opts = append(opts, awscfg.WithEndpointResolverWithOptions(resolver))
	}
	cfg, err := awscfg.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	cli := s3.NewFromConfig(cfg)
	return &S3Uploader{cli: cli, bucket: bucket}, nil
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
