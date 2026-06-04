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

func (r *repository) GetByID(ctx context.Context, id string) (*domain.Availability, error) {
	availabilityID, err := shared.ParseUUID(id)
	if err != nil {
		return nil, err
	}
	row, err := r.queries.GetAvailabilityByID(ctx, availabilityID)
	if err != nil {
		return nil, err
	}
	return &domain.Availability{
		ID:             row.ID.String(),
		ProfessionalID: row.ProfessionalID.String(),
		DayOfWeek:      shared.WeekDay(row.DayOfWeek),
		Shift:          ptrStringToShift(row.Shift),
		StartHour:      row.StartHour,
		EndHour:        row.EndHour,
	}, nil
}

func (r *repository) Create(ctx context.Context, availability *domain.Availability) (*domain.Availability, error) {
	professionalID, err := shared.ParseUUID(availability.ProfessionalID)
	if err != nil {
		return nil, err
	}

	row, err := r.queries.CreateProfessionalAvailability(ctx, CreateProfessionalAvailabilityParams{
		ProfessionalID: professionalID,
		DayOfWeek:      string(availability.DayOfWeek),
		Shift:          shiftToPtr(availability.Shift),
		StartHour:      availability.StartHour,
		EndHour:        availability.EndHour,
	})
	if err != nil {
		return nil, err
	}

	return rowToDomain(row), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	availabilityID, err := shared.ParseUUID(id)
	if err != nil {
		return err
	}
	return r.queries.DeleteProfessionalAvailability(ctx, availabilityID)
}

func (r *repository) DeleteByProfessionalID(ctx context.Context, professionalID string) error {
	id, err := shared.ParseUUID(professionalID)
	if err != nil {
		return err
	}
	return r.queries.DeleteAvailabilitiesByProfessionalID(ctx, id)
}

func (r *repository) GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.Availability, error) {
	professionalIDParsed, err := shared.ParseUUID(professionalID)
	if err != nil {
		return nil, err
	}
	rows, err := r.queries.GetProfessionalAvailabilityByProfessionalID(ctx, professionalIDParsed)
	if err != nil {
		return nil, err
	}
	var availabilities []*domain.Availability
	for _, row := range rows {
		availabilities = append(availabilities, listRowToDomain(row))
	}
	return availabilities, nil
}

func (r *repository) Update(ctx context.Context, availability *domain.Availability) (*domain.Availability, error) {
	availabilityID, err := shared.ParseUUID(availability.ID)
	if err != nil {
		return nil, err
	}

	params := UpdateProfessionalAvailabilityParams{
		ID:        availabilityID,
		StartHour: availability.StartHour,
		EndHour:   availability.EndHour,
	}
	if availability.DayOfWeek != "" {
		d := string(availability.DayOfWeek)
		params.DayOfWeek = &d
	}
	if availability.Shift != nil {
		s := string(*availability.Shift)
		params.Shift = &s
	}

	row, err := r.queries.UpdateProfessionalAvailability(ctx, params)
	if err != nil {
		return nil, err
	}
	return updateRowToDomain(row), nil
}

func rowToDomain(row CreateProfessionalAvailabilityRow) *domain.Availability {
	return &domain.Availability{
		ID:             row.ID.String(),
		ProfessionalID: row.ProfessionalID.String(),
		DayOfWeek:      shared.WeekDay(row.DayOfWeek),
		Shift:          ptrStringToShift(row.Shift),
		StartHour:      row.StartHour,
		EndHour:        row.EndHour,
	}
}

func updateRowToDomain(row UpdateProfessionalAvailabilityRow) *domain.Availability {
	return &domain.Availability{
		ID:             row.ID.String(),
		ProfessionalID: row.ProfessionalID.String(),
		DayOfWeek:      shared.WeekDay(row.DayOfWeek),
		Shift:          ptrStringToShift(row.Shift),
		StartHour:      row.StartHour,
		EndHour:        row.EndHour,
	}
}

func listRowToDomain(row GetProfessionalAvailabilityByProfessionalIDRow) *domain.Availability {
	return &domain.Availability{
		ID:             row.ID.String(),
		ProfessionalID: row.ProfessionalID.String(),
		DayOfWeek:      shared.WeekDay(row.DayOfWeek),
		Shift:          ptrStringToShift(row.Shift),
		StartHour:      row.StartHour,
		EndHour:        row.EndHour,
	}
}

func ptrStringToShift(s *string) *shared.Shift {
	if s == nil {
		return nil
	}
	sh := shared.Shift(*s)
	return &sh
}

func shiftToPtr(s *shared.Shift) *string {
	if s == nil {
		return nil
	}
	str := string(*s)
	return &str
}
