package addresspostgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/diegoHDCz/ajudafio/internal/address/domain"
	"github.com/diegoHDCz/ajudafio/internal/address/ports"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewAddressRepository(db *pgxpool.Pool) ports.AddressRepository {
	return &repository{
		db:      db,
		queries: New(db),
	}
}

func (r *repository) CreateAddress(address *domain.Address) error {
	ctx := context.Background()

	userUID, err := shared.ParseUUID(address.UserID)
	if err != nil {
		return fmt.Errorf("addresspostgres.CreateAddress: invalid userID: %w", err)
	}

	row, err := r.queries.CreateAddress(ctx, CreateAddressParams{
		UserID:      userUID,
		ZipCode:     address.ZipCode,
		AddressLine: address.AddressLine,
		Number:      address.Number,
		Complement:  address.Complement,
		District:    address.District,
		City:        address.City,
		State:       address.State,
		Reference:   address.Reference,
	})
	if err != nil {
		return fmt.Errorf("addresspostgres.CreateAddress: %w", err)
	}

	address.ID = uuid.UUID(row.ID.Bytes).String()
	address.CreatedAt = row.CreatedAt.Time
	address.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (r *repository) GetAddressByID(id string) (*domain.Address, error) {
	ctx := context.Background()

	uid, err := shared.ParseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("addresspostgres.GetAddressByID: invalid id: %w", err)
	}
	row, err := r.queries.GetAddressByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("address not found: %w", err)
		}
		return nil, fmt.Errorf("addresspostgres.GetAddressByID: %w", err)
	}
	return toDomain(row), nil
}

func (r *repository) UpdateAddress(address *domain.Address) error {
	ctx := context.Background()

	uid, err := shared.ParseUUID(address.ID)
	if err != nil {
		return fmt.Errorf("addresspostgres.UpdateAddress: invalid id: %w", err)
	}

	row, err := r.queries.UpdateAddress(ctx, UpdateAddressParams{
		ID:          uid,
		ZipCode:     &address.ZipCode,
		AddressLine: &address.AddressLine,
		Number:      &address.Number,
		Complement:  address.Complement,
		District:    &address.District,
		City:        &address.City,
		State:       &address.State,
		Reference:   address.Reference,
	})
	if err != nil {
		return fmt.Errorf("addresspostgres.UpdateAddress: %w", err)
	}

	address.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (r *repository) DeleteAddress(id string) error {
	ctx := context.Background()

	uid, err := shared.ParseUUID(id)
	if err != nil {
		return fmt.Errorf("addresspostgres.DeleteAddress: invalid id: %w", err)
	}
	if err := r.queries.DeleteAddress(ctx, uid); err != nil {
		return fmt.Errorf("addresspostgres.DeleteAddress: %w", err)
	}
	return nil
}

func (r *repository) GetAddressesByUserID(userID string) ([]*domain.Address, error) {
	ctx := context.Background()

	uid, err := shared.ParseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("addresspostgres.GetAddressesByUserID: invalid userID: %w", err)
	}
	rows, err := r.queries.GetAddressesByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("addresspostgres.GetAddressesByUserID: %w", err)
	}
	result := make([]*domain.Address, 0, len(rows))
	for _, row := range rows {
		result = append(result, toDomain(row))
	}
	return result, nil
}

func (r *repository) GetAllAddresses() ([]*domain.Address, error) {
	ctx := context.Background()

	rows, err := r.queries.GetAllAddresses(ctx, GetAllAddressesParams{})
	if err != nil {
		return nil, fmt.Errorf("addresspostgres.GetAllAddresses: %w", err)
	}
	result := make([]*domain.Address, 0, len(rows))
	for _, row := range rows {
		result = append(result, toDomain(row))
	}
	return result, nil
}

func (r *repository) GetAddressesByCity(city string) ([]*domain.Address, error) {
	ctx := context.Background()

	rows, err := r.queries.GetAddressesByCity(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("addresspostgres.GetAddressesByCity: %w", err)
	}
	result := make([]*domain.Address, 0, len(rows))
	for _, row := range rows {
		result = append(result, toDomain(row))
	}
	return result, nil
}

func toDomain(row Address) *domain.Address {
	return &domain.Address{
		ID:          uuid.UUID(row.ID.Bytes).String(),
		UserID:      uuid.UUID(row.UserID.Bytes).String(),
		ZipCode:     row.ZipCode,
		AddressLine: row.AddressLine,
		Number:      row.Number,
		Complement:  row.Complement,
		District:    row.District,
		City:        row.City,
		State:       row.State,
		Reference:   row.Reference,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}
