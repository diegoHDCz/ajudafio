package appointmentpostgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/diegoHDCz/ajudafio/internal/appointment/domain"
	"github.com/diegoHDCz/ajudafio/internal/appointment/ports"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewRepository(db *pgxpool.Pool) ports.AppointmentRepository {
	return &repository{db: db, queries: New(db)}
}

func (r *repository) CreateAppointment(ctx context.Context, a *domain.Appointment) error {
	contractID, err := shared.ParseUUID(a.ContractID)
	if err != nil {
		return fmt.Errorf("appointmentpostgres.CreateAppointment: invalid contract_id: %w", err)
	}
	clientID, err := shared.ParseUUID(a.ClientID)
	if err != nil {
		return fmt.Errorf("appointmentpostgres.CreateAppointment: invalid client_id: %w", err)
	}
	professionalID, err := shared.ParseUUID(a.ProfessionalID)
	if err != nil {
		return fmt.Errorf("appointmentpostgres.CreateAppointment: invalid professional_id: %w", err)
	}
	row, err := r.queries.CreateAppointment(ctx, CreateAppointmentParams{
		ContractID:     contractID,
		ClientID:       clientID,
		ProfessionalID: professionalID,
		Date:           timeToPgDate(a.Date),
		StartTime:      timeTopgTime(a.StartTime),
		EndTime:        timeTopgTime(a.EndTime),
		ZipCode:        a.ZipCode,
		AddressLine:    a.AddressLine,
		Number:         a.Number,
		Complement:     a.Complement,
		District:       a.District,
		City:           a.City,
		State:          a.State,
		Reference:      a.Reference,
	})
	if err != nil {
		return fmt.Errorf("appointmentpostgres.CreateAppointment: %w", err)
	}
	a.ID = uuid.UUID(row.ID.Bytes).String()
	a.Status = row.Status
	a.Version = int(row.Version)
	a.CreatedAt = row.CreatedAt.Time
	a.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (r *repository) GetAppointmentByID(ctx context.Context, id string) (*domain.Appointment, error) {
	uid, err := shared.ParseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("appointmentpostgres.GetAppointmentByID: invalid id: %w", err)
	}
	row, err := r.queries.GetAppointmentByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("appointment not found: %w", err)
		}
		return nil, fmt.Errorf("appointmentpostgres.GetAppointmentByID: %w", err)
	}
	return rowToDomain(row), nil
}

func (r *repository) GetAppointmentsByContractID(ctx context.Context, contractID string) ([]*domain.Appointment, error) {
	uid, err := shared.ParseUUID(contractID)
	if err != nil {
		return nil, fmt.Errorf("appointmentpostgres.GetAppointmentsByContractID: invalid contract_id: %w", err)
	}
	rows, err := r.queries.GetAppointmentsByContractID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("appointmentpostgres.GetAppointmentsByContractID: %w", err)
	}
	return rowsToDomain(rows), nil
}

func (r *repository) GetAppointmentsByClientID(ctx context.Context, clientID string) ([]*domain.Appointment, error) {
	uid, err := shared.ParseUUID(clientID)
	if err != nil {
		return nil, fmt.Errorf("appointmentpostgres.GetAppointmentsByClientID: invalid client_id: %w", err)
	}
	rows, err := r.queries.GetAppointmentsByClientID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("appointmentpostgres.GetAppointmentsByClientID: %w", err)
	}
	return rowsToDomain(rows), nil
}

func (r *repository) GetAppointmentsByProfessionalID(ctx context.Context, professionalID string) ([]*domain.Appointment, error) {
	uid, err := shared.ParseUUID(professionalID)
	if err != nil {
		return nil, fmt.Errorf("appointmentpostgres.GetAppointmentsByProfessionalID: invalid professional_id: %w", err)
	}
	rows, err := r.queries.GetAppointmentsByProfessionalID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("appointmentpostgres.GetAppointmentsByProfessionalID: %w", err)
	}
	return rowsToDomain(rows), nil
}

func (r *repository) UpdateAppointmentStatus(ctx context.Context, id, status string, version int) (*domain.Appointment, error) {
	uid, err := shared.ParseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("appointmentpostgres.UpdateAppointmentStatus: invalid id: %w", err)
	}
	row, err := r.queries.UpdateAppointmentStatus(ctx, UpdateAppointmentStatusParams{
		ID:      uid,
		Status:  status,
		Version: int32(version),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("appointment not found or version conflict: %w", err)
		}
		return nil, fmt.Errorf("appointmentpostgres.UpdateAppointmentStatus: %w", err)
	}
	return rowToDomain(row), nil
}

func (r *repository) DeleteAppointment(ctx context.Context, id string) error {
	uid, err := shared.ParseUUID(id)
	if err != nil {
		return fmt.Errorf("appointmentpostgres.DeleteAppointment: invalid id: %w", err)
	}
	if err := r.queries.DeleteAppointment(ctx, uid); err != nil {
		return fmt.Errorf("appointmentpostgres.DeleteAppointment: %w", err)
	}
	return nil
}

func (r *repository) HasOverlap(ctx context.Context, professionalID string, date, start, end time.Time) (bool, error) {
	uid, err := shared.ParseUUID(professionalID)
	if err != nil {
		return false, fmt.Errorf("appointmentpostgres.HasOverlap: invalid professional_id: %w", err)
	}
	has, err := r.queries.CheckOverlap(ctx, CheckOverlapParams{
		ProfessionalID: uid,
		Date:           timeToPgDate(date),
		StartTime:      timeTopgTime(start),
		EndTime:        timeTopgTime(end),
	})
	if err != nil {
		return false, fmt.Errorf("appointmentpostgres.HasOverlap: %w", err)
	}
	return has, nil
}

func rowToDomain(row Appointment) *domain.Appointment {
	return &domain.Appointment{
		ID:             uuid.UUID(row.ID.Bytes).String(),
		ContractID:     uuid.UUID(row.ContractID.Bytes).String(),
		ClientID:       uuid.UUID(row.ClientID.Bytes).String(),
		ProfessionalID: uuid.UUID(row.ProfessionalID.Bytes).String(),
		Date:           row.Date.Time,
		StartTime:      pgTimeToTime(row.StartTime),
		EndTime:        pgTimeToTime(row.EndTime),
		Status:         row.Status,
		ZipCode:        row.ZipCode,
		AddressLine:    row.AddressLine,
		Number:         row.Number,
		Complement:     row.Complement,
		District:       row.District,
		City:           row.City,
		State:          row.State,
		Reference:      row.Reference,
		Version:        int(row.Version),
		CreatedAt:      row.CreatedAt.Time,
		UpdatedAt:      row.UpdatedAt.Time,
	}
}

func rowsToDomain(rows []Appointment) []*domain.Appointment {
	result := make([]*domain.Appointment, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToDomain(row))
	}
	return result
}

func pgTimeToTime(t pgtype.Time) time.Time {
	us := t.Microseconds
	h := int(us / (3600 * 1_000_000))
	us %= 3600 * 1_000_000
	m := int(us / (60 * 1_000_000))
	s := int((us % (60 * 1_000_000)) / 1_000_000)
	return time.Date(0, 1, 1, h, m, s, 0, time.UTC)
}

func timeTopgTime(t time.Time) pgtype.Time {
	return pgtype.Time{
		Microseconds: int64(t.Hour()*3600+t.Minute()*60+t.Second()) * 1_000_000,
		Valid:        true,
	}
}

func timeToPgDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}
