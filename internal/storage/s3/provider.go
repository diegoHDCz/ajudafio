package s3provider

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	storagePorts "github.com/diegoHDCz/ajudafio/internal/storage/ports"
)

type Provider struct {
	client *s3.Client
	bucket string
	region string
}

func New(accessKeyID, secretAccessKey, region, bucket string) storagePorts.StorageProvider {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("s3provider: failed to load AWS config: %v", err))
	}
	return &Provider{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
		region: region,
	}
}

func (p *Provider) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	_, err := p.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("s3provider.Upload: %w", err)
	}
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", p.bucket, p.region, key)
	return url, nil
}

func (p *Provider) Delete(ctx context.Context, key string) error {
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("s3provider.Delete: %w", err)
	}
	return nil
}
