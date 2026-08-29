package auditpostgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/diegoHDCz/ajudafio/internal/audit/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewRepository(db *pgxpool.Pool) ports.AuditRepository {
	return &repository{db: db, queries: New(db)}
}

func (r *repository) Insert(ctx context.Context, userID *string, action string, metadata map[string]any) error {
	var uid pgtype.UUID
	if userID != nil {
		parsed, err := uuid.Parse(*userID)
		if err != nil {
			return fmt.Errorf("auditpostgres.Insert: invalid user id: %w", err)
		}
		uid = pgtype.UUID{Bytes: parsed, Valid: true}
	}

	if metadata == nil {
		metadata = map[string]any{}
	}
	raw, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("auditpostgres.Insert: %w", err)
	}

	if err := r.queries.InsertAuditLog(ctx, InsertAuditLogParams{
		UserID:   uid,
		Action:   action,
		Metadata: raw,
	}); err != nil {
		return fmt.Errorf("auditpostgres.Insert: %w", err)
	}
	return nil
}
