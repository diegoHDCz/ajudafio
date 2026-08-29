package ports

import "context"

type AuditRepository interface {
	Insert(ctx context.Context, userID *string, action string, metadata map[string]any) error
}
