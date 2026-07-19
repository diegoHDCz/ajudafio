package authpostgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	"github.com/diegoHDCz/ajudafio/internal/auth/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewRepository(db *pgxpool.Pool) ports.AuthRepository {
	return &repository{
		db:      db,
		queries: New(db),
	}
}

func (r *repository) GetUserAuthByEmail(ctx context.Context, email string) (*authdomain.UserAuth, error) {
	row, err := r.queries.GetUserAuthByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("authpostgres.GetUserAuthByEmail: %w", err)
	}
	return &authdomain.UserAuth{
		ID:           uuid.UUID(row.ID.Bytes).String(),
		Email:        row.Email,
		Name:         row.Name,
		Role:         row.Role,
		PasswordHash: derefString(row.PasswordHash),
	}, nil
}

func (r *repository) SetPasswordHash(ctx context.Context, userID string, passwordHash string) error {
	uid, err := parseUUID(userID)
	if err != nil {
		return fmt.Errorf("authpostgres.SetPasswordHash: invalid id: %w", err)
	}
	if err := r.queries.SetUserPasswordHash(ctx, SetUserPasswordHashParams{
		ID:           uid,
		PasswordHash: &passwordHash,
	}); err != nil {
		return fmt.Errorf("authpostgres.SetPasswordHash: %w", err)
	}
	return nil
}

func (r *repository) CreateRefreshToken(ctx context.Context, token *authdomain.RefreshToken) (*authdomain.RefreshToken, error) {
	uid, err := parseUUID(token.ID)
	if err != nil {
		return nil, fmt.Errorf("authpostgres.CreateRefreshToken: invalid id: %w", err)
	}
	userID, err := parseUUID(token.UserID)
	if err != nil {
		return nil, fmt.Errorf("authpostgres.CreateRefreshToken: invalid user id: %w", err)
	}
	row, err := r.queries.CreateRefreshToken(ctx, CreateRefreshTokenParams{
		ID:        uid,
		UserID:    userID,
		TokenHash: token.TokenHash,
		ExpiresAt: pgtype.Timestamp{Time: token.ExpiresAt, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("authpostgres.CreateRefreshToken: %w", err)
	}
	return mapRefreshToken(row), nil
}

func (r *repository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*authdomain.RefreshToken, error) {
	row, err := r.queries.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("refresh token not found: %w", err)
		}
		return nil, fmt.Errorf("authpostgres.GetRefreshTokenByHash: %w", err)
	}
	return mapRefreshToken(row), nil
}

func (r *repository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	if err := r.queries.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return fmt.Errorf("authpostgres.RevokeRefreshToken: %w", err)
	}
	return nil
}

func (r *repository) RevokeAllRefreshTokensForUser(ctx context.Context, userID string) error {
	uid, err := parseUUID(userID)
	if err != nil {
		return fmt.Errorf("authpostgres.RevokeAllRefreshTokensForUser: invalid user id: %w", err)
	}
	if err := r.queries.RevokeAllRefreshTokensForUser(ctx, uid); err != nil {
		return fmt.Errorf("authpostgres.RevokeAllRefreshTokensForUser: %w", err)
	}
	return nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func mapRefreshToken(row RefreshToken) *authdomain.RefreshToken {
	var revokedAt *time.Time
	if row.RevokedAt.Valid {
		revokedAt = &row.RevokedAt.Time
	}
	return &authdomain.RefreshToken{
		ID:        uuid.UUID(row.ID.Bytes).String(),
		UserID:    uuid.UUID(row.UserID.Bytes).String(),
		TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt.Time,
		RevokedAt: revokedAt,
		CreatedAt: row.CreatedAt.Time,
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func parseUUID(s string) (pgtype.UUID, error) {
	uid, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: uid, Valid: true}, nil
}
