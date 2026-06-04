package ports

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/availability/domain"
)

type AvailabilityRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Availability, error)
	GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.Availability, error)
	Create(ctx context.Context, availability *domain.Availability) (*domain.Availability, error)
	Update(ctx context.Context, availability *domain.Availability) (*domain.Availability, error)
	Delete(ctx context.Context, id string) error
	DeleteByProfessionalID(ctx context.Context, professionalID string) error
}
