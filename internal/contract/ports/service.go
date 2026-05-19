package ports

import (
	"context"
	"time"

	"github.com/diegoHDCz/ajudafio/internal/contract/domain"
	"github.com/diegoHDCz/ajudafio/internal/shared"
)

type ContractService interface {
	Create(ctx context.Context, input CreateContractInput) (*domain.Contract, error)
	GetByID(ctx context.Context, id string) (*domain.Contract, error)
	Update(ctx context.Context, input UpdateContractInput) (*domain.Contract, error)
	Delete(ctx context.Context, id string) error
	GetByUserID(ctx context.Context, userID string) ([]*domain.Contract, error)
	GetAll(ctx context.Context) ([]*domain.Contract, error)
	GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.Contract, error)
	GetByStatus(ctx context.Context, status string) ([]*domain.Contract, error)
	GetByUserIDAndStatus(ctx context.Context, userID string, status string) ([]*domain.Contract, error)
	GetByProfessionalIDAndStatus(ctx context.Context, professionalID string, status string) ([]*domain.Contract, error)
	GetByProfessionalCategory(ctx context.Context, category string) ([]*domain.Contract, error)
}

type CreateContractInput struct {
	ClientID       string
	ProfessionalID string
	HourRate       int
	TotalAmount    int
	Details        []byte
	WeekDays       []shared.WeekDay
	Shift          shared.Shift
	StartTime      time.Time
	HoursPerDay    int
	TotalHours     int
}

type UpdateContractInput struct {
	ID          string
	Status      *string
	HourRate    *int
	TotalAmount *int
	Details     []byte
	WeekDays    []shared.WeekDay
	Shift       *shared.Shift
	StartTime   *time.Time
	HoursPerDay *int
	TotalHours  *int
}
