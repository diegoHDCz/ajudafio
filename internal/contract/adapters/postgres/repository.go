package contractpostgres

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/contract/domain"
	"github.com/diegoHDCz/ajudafio/internal/contract/ports"
	"github.com/diegoHDCz/ajudafio/internal/shared"
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

// DeleteContract implements [ports.ContractRepository].
func (r *repository) DeleteContract(context.Context, string) error {
	panic("unimplemented")
}

// GetAAllContractsByProfessionalCategory implements [ports.ContractRepository].
func (r *repository) GetAAllContractsByProfessionalCategory(context.Context, string) ([]*domain.Contract, error) {
	panic("unimplemented")
}

// GetAAllContractsByProfessionalIDAndStatus implements [ports.ContractRepository].
func (r *repository) GetAAllContractsByProfessionalIDAndStatus(context.Context, string, string) ([]*domain.Contract, error) {
	panic("unimplemented")
}

// GetAllContracts implements [ports.ContractRepository].
func (r *repository) GetAllContracts(context.Context) ([]*domain.Contract, error) {
	panic("unimplemented")
}

// GetAllContractsByProfessionalID implements [ports.ContractRepository].
func (r *repository) GetAllContractsByProfessionalID(context.Context, string) ([]*domain.Contract, error) {
	panic("unimplemented")
}

// GetAllContractsByStatus implements [ports.ContractRepository].
func (r *repository) GetAllContractsByStatus(context.Context, string) ([]*domain.Contract, error) {
	panic("unimplemented")
}

// GetAllContractsByUserIDAndStatus implements [ports.ContractRepository].
func (r *repository) GetAllContractsByUserIDAndStatus(context.Context, string, string) ([]*domain.Contract, error) {
	panic("unimplemented")
}

// GetContractByID implements [ports.ContractRepository].
func (r *repository) GetContractByID(context.Context, string) (*domain.Contract, error) {
	panic("unimplemented")
}

// GetContractsByUserID implements [ports.ContractRepository].
func (r *repository) GetContractsByUserID(context.Context, string) ([]*domain.Contract, error) {
	panic("unimplemented")
}

// UpdateContract implements [ports.ContractRepository].
func (r *repository) UpdateContract(context.Context, *domain.Contract) error {
	panic("unimplemented")
}

func (r *repository) CreateContract(ctx context.Context, contract *domain.Contract) error {
	client, err := shared.ParseUUID(contract.ClientID)
	if err != nil {
		return err
	}
	professional, err := shared.ParseUUID(contract.ProfessionalID)
	if err != nil {
		return err
	}

	daysStr, err := shared.ParseDayOfWeek(contract.WeekDays)
	if err != nil {
		return err
	}
	startTime := pgtype.Time{
		Microseconds: int64(contract.StartTime.Hour()*3600+contract.StartTime.Minute()*60+contract.StartTime.Second()) * 1000000,
		Valid:        true, // O "Valid: true" substitui o antigo "Status: pgtype.Present"
	}
	_, err = r.queries.CreateContract(ctx, CreateContractParams{
		ClientID:       client,
		ProfessionalID: professional,
		HourRate:       int32(contract.HourRate),
		TotalAmount:    int32(contract.TotalAmount),
		Details:        contract.Details,
		WeekDays:       daysStr,
		Shift:          (*string)(&contract.Shift),
		StartTime:      startTime,
		HoursPerDay:    int32(contract.HoursPerDay),
		TotalHours:     int32(contract.TotalHours),
	})
	if err != nil {
		return err
	}
	return nil
}
