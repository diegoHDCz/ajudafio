package profilepostgres

import (
	"context"
	"fmt"

	"github.com/diegoHDCz/ajudafio/internal/profile/domain"
	"github.com/diegoHDCz/ajudafio/internal/profile/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewRepository(db *pgxpool.Pool) ports.ProfileRepository {
	return &repository{db: db, queries: New(db)}
}

func (r *repository) CreateFamilyProfile(ctx context.Context, userID string) (*domain.FamilyProfile, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("profilepostgres.CreateFamilyProfile: invalid user id: %w", err)
	}
	row, err := r.queries.CreateFamilyProfile(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("profilepostgres.CreateFamilyProfile: %w", err)
	}
	return &domain.FamilyProfile{
		ID:        uuid.UUID(row.ID.Bytes).String(),
		UserID:    uuid.UUID(row.UserID.Bytes).String(),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func (r *repository) GetFamilyProfileByUserID(ctx context.Context, userID string) (*domain.FamilyProfile, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("profilepostgres.GetFamilyProfileByUserID: invalid user id: %w", err)
	}
	row, err := r.queries.GetFamilyProfileByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("profilepostgres.GetFamilyProfileByUserID: %w", err)
	}
	return &domain.FamilyProfile{
		ID:        uuid.UUID(row.ID.Bytes).String(),
		UserID:    uuid.UUID(row.UserID.Bytes).String(),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func (r *repository) CreateFinancialProfile(ctx context.Context, userID string) (*domain.FinancialProfile, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("profilepostgres.CreateFinancialProfile: invalid user id: %w", err)
	}
	row, err := r.queries.CreateFinancialProfile(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("profilepostgres.CreateFinancialProfile: %w", err)
	}
	return &domain.FinancialProfile{
		ID:        uuid.UUID(row.ID.Bytes).String(),
		UserID:    uuid.UUID(row.UserID.Bytes).String(),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func (r *repository) GetFinancialProfileByUserID(ctx context.Context, userID string) (*domain.FinancialProfile, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("profilepostgres.GetFinancialProfileByUserID: invalid user id: %w", err)
	}
	row, err := r.queries.GetFinancialProfileByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("profilepostgres.GetFinancialProfileByUserID: %w", err)
	}
	return &domain.FinancialProfile{
		ID:        uuid.UUID(row.ID.Bytes).String(),
		UserID:    uuid.UUID(row.UserID.Bytes).String(),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func parseUUID(s string) (pgtype.UUID, error) {
	uid, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: uid, Valid: true}, nil
}
