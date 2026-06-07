package professionalpostgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/diegoHDCz/ajudafio/internal/professional/domain"
	"github.com/diegoHDCz/ajudafio/internal/professional/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewRepository(db *pgxpool.Pool) ports.ProfessionalRepository {
	return &repository{
		db:      db,
		queries: New(db),
	}
}

func (r *repository) GetByID(ctx context.Context, id string) (*domain.Professional, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("professionalpostgres.GetByID: invalid id: %w", err)
	}
	row, err := r.queries.GetProfessionalByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("professional not found: %w", err)
		}
		return nil, fmt.Errorf("professionalpostgres.GetByID: %w", err)
	}

	p := toDomain(Professional{
		ID:                row.ID,
		UserID:            row.UserID,
		LicenseNumber:     row.LicenseNumber,
		Category:          row.Category,
		YearsOfExperience: row.YearsOfExperience,
		Verified:          row.Verified,
		Resume:            row.Resume,
		Metadata:          row.Metadata,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	})
	p.UserName = &row.UserName
	p.UserAvatarURL = row.UserAvatarUrl
	p.UserEmail = &row.UserEmail
	p.UserRole = &row.UserRole
	return p, nil
}

func (r *repository) GetByUserID(ctx context.Context, userID string) (*domain.Professional, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("professionalpostgres.GetByUserID: invalid userID: %w", err)
	}
	row, err := r.queries.GetProfessionalByUserID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("professional not found: %w", err)
		}
		return nil, fmt.Errorf("professionalpostgres.GetByUserID: %w", err)
	}
	return toDomain(row), nil
}

func (r *repository) Create(ctx context.Context, p *domain.Professional) (*domain.Professional, error) {
	userUID, err := parseUUID(p.UserID)
	if err != nil {
		return nil, fmt.Errorf("professionalpostgres.Create: invalid userID: %w", err)
	}
	yearsOfExp := int32(p.YearsOfExperience)
	row, err := r.queries.CreateProfessional(ctx, CreateProfessionalParams{
		UserID:            userUID,
		LicenseNumber:     &p.LicenseNumber,
		Category:          string(p.Category),
		YearsOfExperience: &yearsOfExp,
		Resume:            &p.Resume,
		Metadata:          p.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("professionalpostgres.Create: %w", err)
	}
	return toDomain(row), nil
}

func (r *repository) Update(ctx context.Context, p *domain.Professional) (*domain.Professional, error) {
	uid, err := parseUUID(p.ID)
	if err != nil {
		return nil, fmt.Errorf("professionalpostgres.Update: invalid id: %w", err)
	}
	category := string(p.Category)
	yearsOfExp := int32(p.YearsOfExperience)
	row, err := r.queries.UpdateProfessional(ctx, UpdateProfessionalParams{
		ID:                uid,
		LicenseNumber:     &p.LicenseNumber,
		Category:          &category,
		YearsOfExperience: &yearsOfExp,
		Verified:          &p.Verified,
		Resume:            &p.Resume,
		Metadata:          p.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("professionalpostgres.Update: %w", err)
	}
	return toDomain(row), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	uid, err := parseUUID(id)
	if err != nil {
		return fmt.Errorf("professionalpostgres.Delete: invalid id: %w", err)
	}
	if err := r.queries.DeleteProfessional(ctx, uid); err != nil {
		return fmt.Errorf("professionalpostgres.Delete: %w", err)
	}
	return nil
}

func (r *repository) FindWithFilters(ctx context.Context, filters ports.ProfessionalFilters) ([]*domain.Professional, int64, error) {
	rows, err := r.queries.ListProfessionals(ctx, ListProfessionalsParams{
		City:      filters.City,
		State:     filters.State,
		DayOfWeek: filters.DayOfWeek,
		Shift:     filters.Shift,
		LimitVal:  int32(filters.PageSize),
		OffsetVal: int32((filters.Page - 1) * filters.PageSize),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("professionalpostgres.FindWithFilters: %w", err)
	}

	total, err := r.queries.CountProfessionals(ctx, CountProfessionalsParams{
		City:      filters.City,
		State:     filters.State,
		DayOfWeek: filters.DayOfWeek,
		Shift:     filters.Shift,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("professionalpostgres.FindWithFilters: %w", err)
	}

	result := make([]*domain.Professional, 0, len(rows))
	for _, row := range rows {
		p := toDomain(Professional{
			ID:                row.ID,
			UserID:            row.UserID,
			LicenseNumber:     row.LicenseNumber,
			Category:          row.Category,
			YearsOfExperience: row.YearsOfExperience,
			Verified:          row.Verified,
			Resume:            row.Resume,
			Metadata:          row.Metadata,
			CreatedAt:         row.CreatedAt,
			UpdatedAt:         row.UpdatedAt,
		})
		p.UserName = row.UserName
		p.UserAvatarURL = row.UserAvatarUrl
		p.UserEmail = row.UserEmail
		p.UserRole = row.UserRole
		result = append(result, p)
	}
	return result, total, nil
}

func toDomain(row Professional) *domain.Professional {
	p := &domain.Professional{
		ID:        uuid.UUID(row.ID.Bytes).String(),
		UserID:    uuid.UUID(row.UserID.Bytes).String(),
		Category:  domain.Category(row.Category),
		Verified:  row.Verified,
		Metadata:  row.Metadata,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
	if row.LicenseNumber != nil {
		p.LicenseNumber = *row.LicenseNumber
	}
	if row.YearsOfExperience != nil {
		p.YearsOfExperience = int(*row.YearsOfExperience)
	}
	if row.Resume != nil {
		p.Resume = *row.Resume
	}
	return p
}

func parseUUID(s string) (pgtype.UUID, error) {
	uid, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: uid, Valid: true}, nil
}
