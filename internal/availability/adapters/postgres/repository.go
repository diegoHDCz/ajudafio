package availabilitypostgres

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/availability/domain"
	"github.com/diegoHDCz/ajudafio/internal/availability/ports"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewRepository(db *pgxpool.Pool) ports.AvailabilityRepository {
	return &repository{
		db:      db,
		queries: New(db),
	}

}

// Create implements [ports.AvailabilityRepository].
func (r *repository) Create(ctx context.Context, availability *domain.Availability) (*domain.Availability, error) {
	professionalID, err := shared.ParseUUID(availability.ProfessionalID)
	if err != nil {
		return nil, err
	}

	days, err := shared.ParseDayOfWeek(availability.DayOfWeek)

	if err != nil {
		return nil, err
	}

	row, err := r.queries.CreateProfessionalAvailability(ctx, CreateProfessionalAvailabilityParams{
		ProfessionalID: professionalID,
		StartHour:      availability.StartHour,
		EndHour:        availability.EndHour,
		DayOfWeek:      days,
	})

	if err != nil {
		return nil, err
	}

	return toDomain(row), nil
}

// Delete implements [ports.AvailabilityRepository].
func (r *repository) Delete(ctx context.Context, id string) error {
	availabilityID, err := shared.ParseUUID(id)
	if err != nil {
		return err
	}
	r.queries.DeleteProfessionalAvailability(ctx, availabilityID)
	return nil
}

// GetByProfessionalID implements [ports.AvailabilityRepository].
func (r *repository) GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.Availability, error) {
	professionalIDParsed, err := shared.ParseUUID(professionalID)
	if err != nil {
		return nil, err
	}
	row, err := r.queries.GetProfessionalAvailabilityByProfessionalID(ctx, professionalIDParsed)
	if err != nil {
		return nil, err
	}
	var availabilities []*domain.Availability
	for _, r := range row {
		availabilities = append(availabilities, createdAvalibilityToDomain(r))
	}
	return availabilities, nil
}

// Update implements [ports.AvailabilityRepository].
func (r *repository) Update(ctx context.Context, availability *domain.Availability) (*domain.Availability, error) {
	availabilityID, err := shared.ParseUUID(availability.ID)
	if err != nil {
		return nil, err
	}

	days, err := shared.ParseDayOfWeek(availability.DayOfWeek)

	if err != nil {
		return nil, err
	}
	row, err := r.queries.UpdateProfessionalAvailability(ctx, UpdateProfessionalAvailabilityParams{
		DayOfWeek: days,
		StartHour: availability.StartHour,
		EndHour:   availability.EndHour,
		ID:        availabilityID,
	})
	if err != nil {
		return nil, err
	}
	return updateRowtoDomain(row), nil
}

func toDomain(row CreateProfessionalAvailabilityRow) *domain.Availability {
	return &domain.Availability{
		ID:             row.ID.String(),
		ProfessionalID: row.ProfessionalID.String(),
		StartHour:      row.StartHour,
		EndHour:        row.EndHour,
		DayOfWeek:      shared.SliceToDayOfWeek(row.DayOfWeek),
	}
}

func updateRowtoDomain(row UpdateProfessionalAvailabilityRow) *domain.Availability {
	return &domain.Availability{
		ID:             row.ID.String(),
		ProfessionalID: row.ProfessionalID.String(),
		StartHour:      row.StartHour,
		EndHour:        row.EndHour,
		DayOfWeek:      shared.SliceToDayOfWeek(row.DayOfWeek),
	}
}

func createdAvalibilityToDomain(row GetProfessionalAvailabilityByProfessionalIDRow) *domain.Availability {
	return &domain.Availability{
		ID:             row.ID.String(),
		ProfessionalID: row.ProfessionalID.String(),
		StartHour:      row.StartHour,
		EndHour:        row.EndHour,
		DayOfWeek:      shared.SliceToDayOfWeek(row.DayOfWeek),
	}
}
