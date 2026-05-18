package ports

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/professional/domain"
)

type ProfessionalRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Professional, error)
	GetByUserID(ctx context.Context, userID string) (*domain.Professional, error)
	Create(ctx context.Context, professional *domain.Professional) (*domain.Professional, error)
	Update(ctx context.Context, professional *domain.Professional) (*domain.Professional, error)
	Delete(ctx context.Context, id string) error
	FindWithFilters(ctx context.Context, filters ProfessionalFilters) ([]*domain.Professional, error)
}
