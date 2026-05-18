package professional

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/professional/domain"
	"github.com/diegoHDCz/ajudafio/internal/professional/ports"
	"github.com/google/uuid"
)

type ProfessionalService struct {
	ps ports.ProfessionalRepository
}

func NewProfessionalService(ps ports.ProfessionalRepository) *ProfessionalService {
	return &ProfessionalService{ps: ps}
}

func (s *ProfessionalService) GetByID(ctx context.Context, id string) (*domain.Professional, error) {
	professional, err := s.ps.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return professional, nil
}

func (s *ProfessionalService) GetByUserID(ctx context.Context, userID string) (*domain.Professional, error) {
	professional, err := s.ps.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return professional, nil
}

func (s *ProfessionalService) Create(ctx context.Context, input ports.CreateProfessionalInput) (*domain.Professional, error) {
	resume := ""
	if input.Resume != nil {
		resume = *input.Resume
	}
	professional, err := domain.NewProfessional(
		uuid.New().String(),
		input.UserID,
		input.LicenseNumber,
		input.Category,
		input.YearsOfExperience,
		resume,
		input.Metadata,
	)
	if err != nil {
		return nil, err
	}
	return s.ps.Create(ctx, professional)
}

func (s *ProfessionalService) Update(ctx context.Context, input ports.UpdateProfessionalInput) (*domain.Professional, error) {
	professional, err := s.ps.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if err := professional.ApplyUpdate(
		input.LicenseNumber,
		input.Category,
		input.YearsOfExperience,
		input.Verified,
		input.Resume,
		input.Metadata,
	); err != nil {
		return nil, err
	}

	return s.ps.Update(ctx, professional)
}

func (s *ProfessionalService) Delete(ctx context.Context, id string) error {
	return s.ps.Delete(ctx, id)
}

func (s *ProfessionalService) FindWithFilters(ctx context.Context, filters ports.ProfessionalFilters) ([]*domain.Professional, error) {
	return s.ps.FindWithFilters(ctx, filters)
}
