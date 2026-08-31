package ports

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/user/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByAuthUserID(ctx context.Context, authUserID string) (*domain.User, error)
	// FindOrCreateByAuthUserID links authUserID to an existing user matched by
	// email, or creates a new one with defaultRole. Safe under concurrent calls.
	FindOrCreateByAuthUserID(ctx context.Context, authUserID, email, name string, defaultRole domain.Role) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	UpdateUserRole(ctx context.Context, id string, role domain.Role) error
	UpdateAvatar(ctx context.Context, id string, avatarURL *string) (*domain.User, error)
	UpdateOnboardingStatus(ctx context.Context, id string, status domain.OnboardingStatus) error
	ListRoles(ctx context.Context, userID string) ([]domain.Role, error)
	// AddRole is idempotent — adding a role the user already has is a no-op.
	AddRole(ctx context.Context, userID string, role domain.Role) error
}
