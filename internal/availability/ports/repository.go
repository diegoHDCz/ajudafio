package ports

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/availability/domain"
)

type AvailabilityRepository interface {
	GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.Availability, error)
	Create(ctx context.Context, availability *domain.Availability) (*domain.Availability, error)
	Update(ctx context.Context, availability *domain.Availability) (*domain.Availability, error)
	Delete(ctx context.Context, id string) error
}
