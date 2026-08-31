package audit

import (
	"context"
	"log/slog"

	"github.com/diegoHDCz/ajudafio/internal/audit/ports"
)

type service struct {
	repo ports.AuditRepository
}

func NewService(repo ports.AuditRepository) ports.AuditService {
	return &service{repo: repo}
}

func (s *service) Log(ctx context.Context, userID *string, action string, metadata map[string]any) {
	if err := s.repo.Insert(ctx, userID, action, metadata); err != nil {
		slog.Error("audit: failed to record event", "action", action, "error", err)
	}
}
