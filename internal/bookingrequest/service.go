package bookingrequest

import (
	"context"
	"errors"
	"fmt"

	"github.com/diegoHDCz/ajudafio/internal/bookingrequest/domain"
	"github.com/diegoHDCz/ajudafio/internal/bookingrequest/ports"
	"github.com/google/uuid"
)

var _ ports.BookingRequestService = (*BookingRequestService)(nil)

type BookingRequestService struct {
	repo ports.BookingRequestRepository
}

func NewBookingRequestService(repo ports.BookingRequestRepository) *BookingRequestService {
	return &BookingRequestService{repo: repo}
}

func (s *BookingRequestService) Create(ctx context.Context, input ports.CreateBookingRequestInput) (*domain.BookingRequest, error) {
	req := &domain.BookingRequest{
		ID:              uuid.New().String(),
		ClientID:        input.ClientID,
		ProfessionalID:  input.ProfessionalID,
		AddressID:       input.AddressID,
		ProposedValue:   input.ProposedValue,
		ScheduleDetails: input.ScheduleDetails,
	}
	return s.repo.Create(ctx, req)
}

func (s *BookingRequestService) GetByID(ctx context.Context, id string) (*domain.BookingRequest, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BookingRequestService) GetByClientID(ctx context.Context, clientID string) ([]*domain.BookingRequest, error) {
	return s.repo.GetByClientID(ctx, clientID)
}

func (s *BookingRequestService) GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.BookingRequest, error) {
	return s.repo.GetByProfessionalID(ctx, professionalID)
}

func (s *BookingRequestService) UpdateStatus(ctx context.Context, id, status string, rejectionReason *string) (*domain.BookingRequest, error) {
	if status == string(domain.StatusRejected) && rejectionReason == nil {
		return nil, errors.New("rejection_reason é obrigatório quando status é REJECTED")
	}
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return nil, fmt.Errorf("bookingrequest.UpdateStatus: %w", err)
	}
	return s.repo.UpdateStatus(ctx, id, status, rejectionReason)
}
