package ports

import (
	"context"

	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	"github.com/google/uuid"
)

type AuthRepository interface {
	GetUserAuthByEmail(ctx context.Context, email string) (*authdomain.UserAuth, error)
	SetPasswordHash(ctx context.Context, userID string, passwordHash string) error

	CreateRefreshToken(ctx context.Context, token *authdomain.RefreshToken) (*authdomain.RefreshToken, error)
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*authdomain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllRefreshTokensForUser(ctx context.Context, userID string) error
	DeleteExpiredRefreshTokens(ctx context.Context) error

	// GetIdentityByProvider returns nil, nil if no identity matches.
	GetIdentityByProvider(ctx context.Context, provider, providerUserID string) (*authdomain.Identity, error)
	CreateIdentity(ctx context.Context, userID uuid.UUID, provider, providerUserID, email string) (*authdomain.Identity, error)
}
