package ports

import "context"

// AuditService writes audit events. It swallows persistence errors after
// logging them (best-effort — see docs/result/refactor-auth.md open question
// #7): an auditing failure must never block the business operation it
// describes. Security-critical call sites that need a hard failure guarantee
// should call the repository directly within their own transaction instead.
type AuditService interface {
	Log(ctx context.Context, userID *string, action string, metadata map[string]any)
}
