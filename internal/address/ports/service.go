package ports

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/address/domain"
)

type AddressService interface {
	GetByID(ctx context.Context, id string) (*domain.Address, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Address, error)
	GetByContractID(ctx context.Context, contractID string) ([]*domain.Address, error)
	Create(ctx context.Context, input CreateAddressInput) (*domain.Address, error)
	Update(ctx context.Context, input UpdateAddressInput) (*domain.Address, error)
	Delete(ctx context.Context, id string) error
}

type CreateAddressInput struct {
	UserID      string
	ContractID  *string
	ZipCode     string
	AddressLine string
	Number      string
	Complement  *string
	District    string
	City        string
	State       string
	Reference   *string
}

type UpdateAddressInput struct {
	ID          string
	ZipCode     *string
	AddressLine *string
	Number      *string
	Complement  *string
	District    *string
	City        *string
	State       *string
	Reference   *string
}
