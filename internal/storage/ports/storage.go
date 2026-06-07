package ports

import (
	"context"
	"time"
)

type StorageProvider interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
	GetSignedURL(ctx context.Context, key string, expiresIn time.Duration) (string, error)
}
