package ports

import (
	"context"

	"github.com/seuuser/healthcontracts/internal/shared/valueobjects"
	"github.com/seuuser/healthcontracts/internal/users/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id valueobjects.EntityID) (*domain.User, error)
}

type AccountRepository interface {
	SaveCredential(ctx context.Context, acc domain.CredentialAccount) error
	FindByUserID(ctx context.Context, userID valueobjects.EntityID) ([]domain.Account, error)
}

type SessionRepository interface {
	Save(ctx context.Context, session *domain.Session) error
	FindByToken(ctx context.Context, token string) (*domain.Session, error)
	Delete(ctx context.Context, token string) error
}
