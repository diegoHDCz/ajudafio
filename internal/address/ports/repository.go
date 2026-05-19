package ports

import "github.com/diegoHDCz/ajudafio/internal/address/domain"

type AddressRepository interface {
	CreateAddress(address *domain.Address) error
	GetAddressByID(id string) (*domain.Address, error)
	UpdateAddress(address *domain.Address) error
	DeleteAddress(id string) error
	GetAddressesByUserID(userID string) ([]*domain.Address, error)
	GetAddressesByContractID(contractID string) ([]*domain.Address, error)
	GetAllAddresses() ([]*domain.Address, error)
	GetAddressesByCity(city string) ([]*domain.Address, error)
}
