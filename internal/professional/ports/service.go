package ports

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/professional/domain"
)

type ProfessionalService interface {
	GetByID(ctx context.Context, id string) (*domain.Professional, error)
	GetByUserID(ctx context.Context, userID string) (*domain.Professional, error)
	Create(ctx context.Context, input CreateProfessionalInput) (*domain.Professional, error)
	Update(ctx context.Context, input UpdateProfessionalInput) (*domain.Professional, error)
	Delete(ctx context.Context, id string) error
	FindWithFilters(ctx context.Context, filters ProfessionalFilters) ([]*domain.Professional, error)
}

type CreateProfessionalInput struct {
	UserID            string
	LicenseNumber     string
	Category          domain.Category
	YearsOfExperience int
	Resume            *string
	Metadata          []byte
}

type UpdateProfessionalInput struct {
	ID                string
	LicenseNumber     *string
	Category          *domain.Category
	YearsOfExperience *int
	Verified          *bool
	Resume            *string
	Metadata          []byte
}

type ProfessionalFilters struct {
	City      *string
	State     *string
	DayOfWeek []string
	Shift     []string
}
