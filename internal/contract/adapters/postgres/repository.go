package contractpostgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/diegoHDCz/ajudafio/internal/contract/domain"
	"github.com/diegoHDCz/ajudafio/internal/contract/ports"
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

func NewRepository(db *pgxpool.Pool) ports.ContractRepository {
	return &repository{db: db, queries: New(db)}
}

func (r *repository) CreateContract(ctx context.Context, contract *domain.Contract) error {
	client, err := shared.ParseUUID(contract.ClientID)
	if err != nil {
		return fmt.Errorf("contractpostgres.CreateContract: invalid clientID: %w", err)
	}
	professional, err := shared.ParseUUID(contract.ProfessionalID)
	if err != nil {
		return fmt.Errorf("contractpostgres.CreateContract: invalid professionalID: %w", err)
	}
	daysStr, err := shared.ParseDayOfWeek(contract.WeekDays)
	if err != nil {
		return fmt.Errorf("contractpostgres.CreateContract: %w", err)
	}
	shift := string(contract.Shift)
	row, err := r.queries.CreateContract(ctx, CreateContractParams{
		ClientID:       client,
		ProfessionalID: professional,
		HourRate:       int32(contract.HourRate),
		TotalAmount:    int32(contract.TotalAmount),
		Details:        contract.Details,
		WeekDays:       daysStr,
		Shift:          &shift,
		StartTime:      timeTopgTime(contract.StartTime),
		HoursPerDay:    int32(contract.HoursPerDay),
		TotalHours:     int32(contract.TotalHours),
	})
	if err != nil {
		return fmt.Errorf("contractpostgres.CreateContract: %w", err)
	}
	contract.ID = uuid.UUID(row.ID.Bytes).String()
	contract.CreatedAt = row.CreatedAt.Time
	return nil
}

func (r *repository) GetContractByID(ctx context.Context, contractID string) (*domain.Contract, error) {
	id, err := shared.ParseUUID(contractID)
	if err != nil {
		return nil, fmt.Errorf("contractpostgres.GetContractByID: invalid id: %w", err)
	}
	row, err := r.queries.GetContractByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("contract not found: %w", err)
		}
		return nil, fmt.Errorf("contractpostgres.GetContractByID: %w", err)
	}
	return rowToDomain(row.ID, row.ClientID, row.ProfessionalID, row.Status, row.HourRate, row.TotalAmount, row.Details, row.WeekDays, row.Shift, row.StartTime, row.HoursPerDay, row.TotalHours, row.CreatedAt), nil
}

func (r *repository) UpdateContract(ctx context.Context, contract *domain.Contract) error {
	id, err := shared.ParseUUID(contract.ID)
	if err != nil {
		return fmt.Errorf("contractpostgres.UpdateContract: invalid id: %w", err)
	}
	daysStr, err := shared.ParseDayOfWeek(contract.WeekDays)
	if err != nil {
		return fmt.Errorf("contractpostgres.UpdateContract: %w", err)
	}
	status := contract.Status
	hourRate := int32(contract.HourRate)
	totalAmount := int32(contract.TotalAmount)
	hoursPerDay := int32(contract.HoursPerDay)
	totalHours := int32(contract.TotalHours)
	shift := string(contract.Shift)
	_, err = r.queries.UpdateContract(ctx, UpdateContractParams{
		ID:          id,
		Status:      &status,
		HourRate:    &hourRate,
		TotalAmount: &totalAmount,
		Details:     contract.Details,
		WeekDays:    daysStr,
		Shift:       &shift,
		StartTime:   timeTopgTime(contract.StartTime),
		HoursPerDay: &hoursPerDay,
		TotalHours:  &totalHours,
	})
	if err != nil {
		return fmt.Errorf("contractpostgres.UpdateContract: %w", err)
	}
	return nil
}

func (r *repository) DeleteContract(ctx context.Context, contractID string) error {
	id, err := shared.ParseUUID(contractID)
	if err != nil {
		return fmt.Errorf("contractpostgres.DeleteContract: invalid id: %w", err)
	}
	if err := r.queries.DeleteContract(ctx, id); err != nil {
		return fmt.Errorf("contractpostgres.DeleteContract: %w", err)
	}
	return nil
}

func (r *repository) GetContractsByUserID(ctx context.Context, userID string) ([]*domain.Contract, error) {
	id, err := shared.ParseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("contractpostgres.GetContractsByUserID: invalid userID: %w", err)
	}
	rows, err := r.queries.GetContractsByUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("contractpostgres.GetContractsByUserID: %w", err)
	}
	result := make([]*domain.Contract, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToDomain(row.ID, row.ClientID, row.ProfessionalID, row.Status, row.HourRate, row.TotalAmount, row.Details, row.WeekDays, row.Shift, row.StartTime, row.HoursPerDay, row.TotalHours, row.CreatedAt))
	}
	return result, nil
}

func (r *repository) GetAllContracts(ctx context.Context) ([]*domain.Contract, error) {
	rows, err := r.queries.GetAllContracts(ctx)
	if err != nil {
		return nil, fmt.Errorf("contractpostgres.GetAllContracts: %w", err)
	}
	result := make([]*domain.Contract, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToDomain(row.ID, row.ClientID, row.ProfessionalID, row.Status, row.HourRate, row.TotalAmount, row.Details, row.WeekDays, row.Shift, row.StartTime, row.HoursPerDay, row.TotalHours, row.CreatedAt))
	}
	return result, nil
}

func (r *repository) GetAllContractsByProfessionalID(ctx context.Context, professionalID string) ([]*domain.Contract, error) {
	id, err := shared.ParseUUID(professionalID)
	if err != nil {
		return nil, fmt.Errorf("contractpostgres.GetAllContractsByProfessionalID: invalid professionalID: %w", err)
	}
	rows, err := r.queries.GetAllContractsByProfessionalID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("contractpostgres.GetAllContractsByProfessionalID: %w", err)
	}
	result := make([]*domain.Contract, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToDomain(row.ID, row.ClientID, row.ProfessionalID, row.Status, row.HourRate, row.TotalAmount, row.Details, row.WeekDays, row.Shift, row.StartTime, row.HoursPerDay, row.TotalHours, row.CreatedAt))
	}
	return result, nil
}

func (r *repository) GetAllContractsByStatus(ctx context.Context, status string) ([]*domain.Contract, error) {
	rows, err := r.queries.GetAllContractsByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("contractpostgres.GetAllContractsByStatus: %w", err)
	}
	result := make([]*domain.Contract, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToDomain(row.ID, row.ClientID, row.ProfessionalID, row.Status, row.HourRate, row.TotalAmount, row.Details, row.WeekDays, row.Shift, row.StartTime, row.HoursPerDay, row.TotalHours, row.CreatedAt))
	}
	return result, nil
}

func (r *repository) GetAllContractsByUserIDAndStatus(ctx context.Context, userID string, status string) ([]*domain.Contract, error) {
	id, err := shared.ParseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("contractpostgres.GetAllContractsByUserIDAndStatus: invalid userID: %w", err)
	}
	rows, err := r.queries.GetAllContractsByUserIDAndStatus(ctx, GetAllContractsByUserIDAndStatusParams{
		ClientID: id,
		Status:   status,
	})
	if err != nil {
		return nil, fmt.Errorf("contractpostgres.GetAllContractsByUserIDAndStatus: %w", err)
	}
	result := make([]*domain.Contract, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToDomain(row.ID, row.ClientID, row.ProfessionalID, row.Status, row.HourRate, row.TotalAmount, row.Details, row.WeekDays, row.Shift, row.StartTime, row.HoursPerDay, row.TotalHours, row.CreatedAt))
	}
	return result, nil
}

func (r *repository) GetAAllContractsByProfessionalIDAndStatus(ctx context.Context, professionalID string, status string) ([]*domain.Contract, error) {
	id, err := shared.ParseUUID(professionalID)
	if err != nil {
		return nil, fmt.Errorf("contractpostgres.GetAAllContractsByProfessionalIDAndStatus: invalid professionalID: %w", err)
	}
	rows, err := r.queries.GetAllContractsByProfessionalIDAndStatus(ctx, GetAllContractsByProfessionalIDAndStatusParams{
		ProfessionalID: id,
		Status:         status,
	})
	if err != nil {
		return nil, fmt.Errorf("contractpostgres.GetAAllContractsByProfessionalIDAndStatus: %w", err)
	}
	result := make([]*domain.Contract, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToDomain(row.ID, row.ClientID, row.ProfessionalID, row.Status, row.HourRate, row.TotalAmount, row.Details, row.WeekDays, row.Shift, row.StartTime, row.HoursPerDay, row.TotalHours, row.CreatedAt))
	}
	return result, nil
}

func (r *repository) GetAAllContractsByProfessionalCategory(ctx context.Context, category string) ([]*domain.Contract, error) {
	rows, err := r.queries.GetAllContractsByProfessionalCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("contractpostgres.GetAAllContractsByProfessionalCategory: %w", err)
	}
	result := make([]*domain.Contract, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToDomain(row.ID, row.ClientID, row.ProfessionalID, row.Status, row.HourRate, row.TotalAmount, row.Details, row.WeekDays, row.Shift, row.StartTime, row.HoursPerDay, row.TotalHours, row.CreatedAt))
	}
	return result, nil
}

func rowToDomain(
	id, clientID, professionalID pgtype.UUID,
	status string,
	hourRate, totalAmount int32,
	details []byte,
	weekDays []string,
	shift *string,
	startTime pgtype.Time,
	hoursPerDay, totalHours int32,
	createdAt pgtype.Timestamp,
) *domain.Contract {
	c := &domain.Contract{
		ID:             uuid.UUID(id.Bytes).String(),
		ClientID:       uuid.UUID(clientID.Bytes).String(),
		ProfessionalID: uuid.UUID(professionalID.Bytes).String(),
		Status:         status,
		HourRate:       int(hourRate),
		TotalAmount:    int(totalAmount),
		Details:        details,
		WeekDays:       shared.SliceToDayOfWeek(weekDays),
		HoursPerDay:    int(hoursPerDay),
		TotalHours:     int(totalHours),
		CreatedAt:      createdAt.Time,
		StartTime:      pgTimeToTime(startTime),
	}
	if shift != nil {
		c.Shift = shared.Shift(*shift)
	}
	return c
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
