package ports

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/bookingrequest/domain"
)

type BookingRequestService interface {
	Create(ctx context.Context, input CreateBookingRequestInput) (*domain.BookingRequest, error)
	GetByID(ctx context.Context, id string) (*domain.BookingRequest, error)
	GetByClientID(ctx context.Context, clientID string) ([]*domain.BookingRequest, error)
	GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.BookingRequest, error)
	UpdateStatus(ctx context.Context, id, status string, rejectionReason *string) (*domain.BookingRequest, error)
}

type CreateBookingRequestInput struct {
	ClientID        string
	ProfessionalID  string
	AddressID       string
	ProposedValue   float64
	ScheduleDetails domain.ScheduleDetails
}
