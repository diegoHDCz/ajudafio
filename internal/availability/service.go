package availability

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/availability/domain"
	"github.com/diegoHDCz/ajudafio/internal/availability/ports"
)

type AvailabilityService struct {
	ar ports.AvailabilityRepository
}

func NewAvailabilityService(ar ports.AvailabilityRepository) *AvailabilityService {
	return &AvailabilityService{ar: ar}
}

func (s *AvailabilityService) GetByID(ctx context.Context, id string) (*domain.Availability, error) {
	return s.ar.GetByID(ctx, id)
}

func (s *AvailabilityService) GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.Availability, error) {
	return s.ar.GetByProfessionalID(ctx, professionalID)
}

func (s *AvailabilityService) Create(ctx context.Context, availability *domain.Availability) (*domain.Availability, error) {
	return s.ar.Create(ctx, availability)
}

func (s *AvailabilityService) Update(ctx context.Context, availability *domain.Availability) (*domain.Availability, error) {
	return s.ar.Update(ctx, availability)
}
func (s *AvailabilityService) Delete(ctx context.Context, id string) error {
	return s.ar.Delete(ctx, id)
}
