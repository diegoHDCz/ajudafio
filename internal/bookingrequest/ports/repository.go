package ports

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/bookingrequest/domain"
)

type BookingRequestRepository interface {
	Create(ctx context.Context, req *domain.BookingRequest) (*domain.BookingRequest, error)
	GetByID(ctx context.Context, id string) (*domain.BookingRequest, error)
	GetByClientID(ctx context.Context, clientID string) ([]*domain.BookingRequest, error)
	GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.BookingRequest, error)
	UpdateStatus(ctx context.Context, id, status string, rejectionReason *string) (*domain.BookingRequest, error)
}
