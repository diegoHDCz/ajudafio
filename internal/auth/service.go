package auth

import (
	"context"
	"time"

	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	"github.com/diegoHDCz/ajudafio/internal/auth/ports"
	userdomain "github.com/diegoHDCz/ajudafio/internal/user/domain"
)

type AuthService struct {
	repo ports.AuthRepository
}

func NewAuthService(repo ports.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

// FindOrCreateUser resolves the local user for the given Keycloak claims,
// bootstrapping both the user row and the Keycloak account link on first login.
func (s *AuthService) FindOrCreateUser(ctx context.Context, claims authdomain.Claims) (*userdomain.User, error) {
	account, err := s.repo.GetAccountByProvider(ctx, "keycloak", claims.UserID)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.FindOrCreateUser(ctx, claims)
	if err != nil {
		return nil, err
	}

	if account == nil {
		now := time.Now()
		_, err = s.repo.CreateAccount(ctx, ports.CreateAccountParams{
			AccountID:             claims.UserID,
			ProviderID:            "keycloak",
			UserID:                string(user.ID),
			AccessTokenExpiresAt:  now,
			RefreshTokenExpiresAt: now,
		})
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

// GetAccountsByUser returns all provider accounts linked to the given user ID.
func (s *AuthService) GetAccountsByUser(ctx context.Context, userID string) ([]ports.Account, error) {
	return s.repo.GetAccountsByUserID(ctx, userID)
}
