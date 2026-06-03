package ports

import (
	"context"
	"time"

	"github.com/diegoHDCz/ajudafio/internal/appointment/domain"
)

type AppointmentRepository interface {
	CreateAppointment(ctx context.Context, a *domain.Appointment) error
	GetAppointmentByID(ctx context.Context, id string) (*domain.Appointment, error)
	GetAppointmentsByContractID(ctx context.Context, contractID string) ([]*domain.Appointment, error)
	GetAppointmentsByClientID(ctx context.Context, clientID string) ([]*domain.Appointment, error)
	GetAppointmentsByProfessionalID(ctx context.Context, professionalID string) ([]*domain.Appointment, error)
	UpdateAppointmentStatus(ctx context.Context, id, status string, version int) (*domain.Appointment, error)
	DeleteAppointment(ctx context.Context, id string) error
	HasOverlap(ctx context.Context, professionalID string, date, start, end time.Time) (bool, error)
}
