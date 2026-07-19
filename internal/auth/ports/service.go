package ports

import (
	"context"

	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*authdomain.TokenPair, error)
	Login(ctx context.Context, email, password string) (*authdomain.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*authdomain.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}

type RegisterInput struct {
	Name     string
	Email    string
	Phone    *string
	Password string
}
