package ports

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/user/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id domain.UserID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id domain.UserID) error
}