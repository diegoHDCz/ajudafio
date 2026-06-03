package ports

import (
	"context"
	"time"

	"github.com/diegoHDCz/ajudafio/internal/appointment/domain"
)

type AppointmentService interface {
	Create(ctx context.Context, input CreateAppointmentInput) (*domain.Appointment, error)
	GetByID(ctx context.Context, id string) (*domain.Appointment, error)
	GetByContractID(ctx context.Context, contractID string) ([]*domain.Appointment, error)
	GetByClientID(ctx context.Context, clientID string) ([]*domain.Appointment, error)
	GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.Appointment, error)
	UpdateStatus(ctx context.Context, id, status string) (*domain.Appointment, error)
	Delete(ctx context.Context, id string) error
}

type CreateAppointmentInput struct {
	ContractID     string
	ClientID       string
	ProfessionalID string
	Date           time.Time
	StartTime      time.Time
	EndTime        time.Time
	ZipCode        string
	AddressLine    string
	Number         string
	Complement     *string
	District       string
	City           string
	State          string
	Reference      *string
}
