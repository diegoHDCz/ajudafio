package userpostgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/diegoHDCz/ajudafio/internal/user/domain"
	"github.com/diegoHDCz/ajudafio/internal/user/ports"
	"github.com/jackc/pgx/v5"
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

func (r *repository) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("userpostgres.GetByID: %w", err)
	}
	return toDomain(row), nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("userpostgres.GetByEmail: %w", err)
	}
	return toDomain(row), nil
}

func (r *repository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	row, err := r.queries.CreateUser(ctx, CreateUserParams{
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone, // sqlc gerou como *string, certifique-se que o domínio também seja ou converta aqui
		Role:  string(user.Role),
	})
	if err != nil {
		return nil, fmt.Errorf("userpostgres.Create: %w", err)
	}
	return toDomain(row), nil
}

func (r *repository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	row, err := r.queries.UpdateUser(ctx, UpdateUserParams{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Role:  string(user.Role),
	})
	if err != nil {
		return nil, fmt.Errorf("userpostgres.Update: %w", err)
	}
	return toDomain(row), nil
}

func (r *repository) Delete(ctx context.Context, id domain.UserID) error {
	if err := r.queries.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("userpostgres.Delete: %w", err)
	}
	return nil
}

// toDomain mapeia a struct local "User" (gerada pelo sqlc) para a entidade de domínio.
func toDomain(row User) *domain.User {
	return &domain.User{
		ID:        row.ID, // Já mapeado via sqlc override para domain.UserID
		Name:      row.Name,
		Email:     row.Email,
		Phone:     row.Phone, // *string
		Role:      domain.Role(row.Role),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
