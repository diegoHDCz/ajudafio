package address

import (
	"context"
	"fmt"

	"github.com/diegoHDCz/ajudafio/internal/address/domain"
	"github.com/diegoHDCz/ajudafio/internal/address/ports"
	"github.com/google/uuid"
)

type AddressService struct {
	repo ports.AddressRepository
}

func NewAddressService(repo ports.AddressRepository) *AddressService {
	return &AddressService{repo: repo}
}

func (s *AddressService) GetByID(_ context.Context, id string) (*domain.Address, error) {
	return s.repo.GetAddressByID(id)
}

func (s *AddressService) GetByUserID(_ context.Context, userID string) ([]*domain.Address, error) {
	return s.repo.GetAddressesByUserID(userID)
}

func (s *AddressService) GetByContractID(_ context.Context, contractID string) ([]*domain.Address, error) {
	return s.repo.GetAddressesByContractID(contractID)
}

func (s *AddressService) Create(_ context.Context, input ports.CreateAddressInput) (*domain.Address, error) {
	address := &domain.Address{
		ID:          uuid.New().String(),
		UserID:      input.UserID,
		ContractID:  input.ContractID,
		ZipCode:     input.ZipCode,
		AddressLine: input.AddressLine,
		Number:      input.Number,
		Complement:  input.Complement,
		District:    input.District,
		City:        input.City,
		State:       input.State,
		Reference:   input.Reference,
	}
	if err := s.repo.CreateAddress(address); err != nil {
		return nil, fmt.Errorf("address.Create: %w", err)
	}
	return address, nil
}

func (s *AddressService) Update(_ context.Context, input ports.UpdateAddressInput) (*domain.Address, error) {
	address, err := s.repo.GetAddressByID(input.ID)
	if err != nil {
		return nil, fmt.Errorf("address.Update: %w", err)
	}

	if input.ZipCode != nil {
		address.ZipCode = *input.ZipCode
	}
	if input.AddressLine != nil {
		address.AddressLine = *input.AddressLine
	}
	if input.Number != nil {
		address.Number = *input.Number
	}
	if input.Complement != nil {
		address.Complement = input.Complement
	}
	if input.District != nil {
		address.District = *input.District
	}
	if input.City != nil {
		address.City = *input.City
	}
	if input.State != nil {
		address.State = *input.State
	}
	if input.Reference != nil {
		address.Reference = input.Reference
	}

	if err := s.repo.UpdateAddress(address); err != nil {
		return nil, fmt.Errorf("address.Update: %w", err)
	}
	return address, nil
}

func (s *AddressService) Delete(_ context.Context, id string) error {
	return s.repo.DeleteAddress(id)
}
