package ports

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/user/domain"
)

type UserService interface {
	GetByID(ctx context.Context, id domain.UserID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, input CreateUserInput) (*domain.User, error)
	Update(ctx context.Context, input UpdateUserInput) (*domain.User, error)
	Delete(ctx context.Context, id domain.UserID) error
}

type CreateUserInput struct {
	Email string
	Name  string  // Obrigatório (NOT NULL no banco)
	Phone *string // Opcional no banco
	Role  domain.Role
}

type UpdateUserInput struct {
	ID    domain.UserID
	Email *string
	Name  *string
	Phone *string
	Role  *domain.Role
}
