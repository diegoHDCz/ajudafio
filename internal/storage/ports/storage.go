package ports

import "context"

type StorageProvider interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
}
