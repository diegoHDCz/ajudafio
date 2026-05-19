package ports

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/availability/domain"
)

type AvalabilityService interface {
	GetByID(ctx context.Context, id string) (*domain.Availability, error)
	GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.Availability, error)
	Create(ctx context.Context, input CreateAvailabilityInput) (*domain.Availability, error)
	Update(ctx context.Context, input UpdateAvailabilityInput) (*domain.Availability, error)
	Delete(ctx context.Context, id string) error
}

type CreateAvailabilityInput struct {
	ProfessionalID string
	StartTime      string
	EndTime        string
	Metadata       map[string]interface{}
}

type UpdateAvailabilityInput struct {
	ID             string
	ProfessionalID string
	StartTime      string
	EndTime        string
	Metadata       map[string]interface{}
}
