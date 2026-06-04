package userpostgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/diegoHDCz/ajudafio/internal/user/domain"
	"github.com/diegoHDCz/ajudafio/internal/user/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewRepository(db *pgxpool.Pool) ports.UserRepository {
	return &repository{
		db:      db,
		queries: New(db),
	}
}

func (r *repository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("userpostgres.GetByID: invalid id: %w", err)
	}
	row, err := r.queries.GetUserByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("userpostgres.GetByID: %w", err)
	}
	return mapUser(row.ID, row.Name, row.Email, row.Phone, row.Role, row.AvatarUrl, row.CreatedAt, row.UpdatedAt), nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("userpostgres.GetByEmail: %w", err)
	}
	return mapUser(row.ID, row.Name, row.Email, row.Phone, row.Role, row.AvatarUrl, row.CreatedAt, row.UpdatedAt), nil
}

func (r *repository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	uid, err := parseUUID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("userpostgres.Create: invalid id: %w", err)
	}
	row, err := r.queries.CreateUser(ctx, CreateUserParams{
		ID:    uid,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Role:  string(user.Role),
	})
	if err != nil {
		return nil, fmt.Errorf("userpostgres.Create: %w", err)
	}
	return mapUser(row.ID, row.Name, row.Email, row.Phone, row.Role, row.AvatarUrl, row.CreatedAt, row.UpdatedAt), nil
}

func (r *repository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	uid, err := parseUUID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("userpostgres.Update: invalid id: %w", err)
	}
	row, err := r.queries.UpdateUser(ctx, UpdateUserParams{
		ID:    uid,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Role:  string(user.Role),
	})
	if err != nil {
		return nil, fmt.Errorf("userpostgres.Update: %w", err)
	}
	return mapUser(row.ID, row.Name, row.Email, row.Phone, row.Role, row.AvatarUrl, row.CreatedAt, row.UpdatedAt), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	uid, err := parseUUID(id)
	if err != nil {
		return fmt.Errorf("userpostgres.Delete: invalid id: %w", err)
	}
	if err := r.queries.DeleteUser(ctx, uid); err != nil {
		return fmt.Errorf("userpostgres.Delete: %w", err)
	}
	return nil
}

func (r *repository) UpdateUserRole(ctx context.Context, id string, role domain.Role) error {
	uid, err := parseUUID(id)
	if err != nil {
		return fmt.Errorf("userpostgres.UpdateUserRole: invalid id: %w", err)
	}
	if err := r.queries.UpdateUserRole(ctx, UpdateUserRoleParams{
		ID:   uid,
		Role: string(role),
	}); err != nil {
		return fmt.Errorf("userpostgres.UpdateUserRole: %w", err)
	}
	return nil
}

func (r *repository) UpdateAvatar(ctx context.Context, id string, avatarURL *string) (*domain.User, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("userpostgres.UpdateAvatar: invalid id: %w", err)
	}
	row, err := r.queries.UpdateUserAvatar(ctx, UpdateUserAvatarParams{
		ID:        uid,
		AvatarUrl: avatarURL,
	})
	if err != nil {
		return nil, fmt.Errorf("userpostgres.UpdateAvatar: %w", err)
	}
	return mapUser(row.ID, row.Name, row.Email, row.Phone, row.Role, row.AvatarUrl, row.CreatedAt, row.UpdatedAt), nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func mapUser(
	id pgtype.UUID,
	name, email string,
	phone *string,
	role string,
	avatarURL *string,
	createdAt, updatedAt pgtype.Timestamp,
) *domain.User {
	return &domain.User{
		ID:        uuid.UUID(id.Bytes).String(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Role:      domain.Role(role),
		AvatarURL: avatarURL,
		CreatedAt: createdAt.Time,
		UpdatedAt: updatedAt.Time,
	}
}

func parseUUID(s string) (pgtype.UUID, error) {
	uid, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: uid, Valid: true}, nil
}
