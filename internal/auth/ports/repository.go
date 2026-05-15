package ports

import (
	"context"
	"time"

	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	userdomain "github.com/diegoHDCz/ajudafio/internal/user/domain"
)


type CreateAccountParams struct {
	AccountID             string // Keycloak subject (sub claim)
	ProviderID            string // e.g. "keycloak"
	UserID                string
	AccessToken           *string
	RefreshToken          *string
	IDToken               *string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	Scope                 *string
}


type Account struct {
	ID                    string
	AccountID             string
	ProviderID            string
	UserID                string
	AccessToken           *string
	RefreshToken          *string
	IDToken               *string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	Scope                 *string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}


type CreateSessionParams struct {
	ExpiresAt time.Time
	Token     string
	IPAddress *string
	UserAgent *string
	UserID    string
}


type Session struct {
	ID        string
	ExpiresAt time.Time
	Token     string
	IPAddress *string
	UserAgent *string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateVerificationParams struct {
	Identifier string
	Value      string
	ExpiresAt  time.Time
}

type Verification struct {
	ID         string
	Identifier string
	Value      string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
type AuthRepository interface {

	FindOrCreateUser(ctx context.Context, claims authdomain.Claims) (*userdomain.User, error)

	GetAccountByProvider(ctx context.Context, providerID, accountID string) (*Account, error)

	CreateAccount(ctx context.Context, params CreateAccountParams) (*Account, error)

	GetAccountsByUserID(ctx context.Context, userID string) ([]Account, error)

	CreateSession(ctx context.Context, params CreateSessionParams) (*Session, error)

	GetSessionByToken(ctx context.Context, token string) (*Session, error)

	DeleteSession(ctx context.Context, token string) error

	DeleteExpiredSessions(ctx context.Context) error

	CreateVerification(ctx context.Context, params CreateVerificationParams) (*Verification, error)

	GetVerification(ctx context.Context, identifier, value string) (*Verification, error)

	DeleteVerification(ctx context.Context, identifier, value string) error
}
