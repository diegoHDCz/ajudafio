package contract

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/contract/domain"
	"github.com/diegoHDCz/ajudafio/internal/contract/ports"
	"github.com/google/uuid"
)

var _ ports.ContractService = (*ContractService)(nil)

type ContractService struct {
	repo ports.ContractRepository
}

func NewContractService(repo ports.ContractRepository) *ContractService {
	return &ContractService{repo: repo}
}

func (s *ContractService) Create(ctx context.Context, input ports.CreateContractInput) (*domain.Contract, error) {
	contract := &domain.Contract{
		ID:             uuid.New().String(),
		ClientID:       input.ClientID,
		ProfessionalID: input.ProfessionalID,
		HourRate:       input.HourRate,
		TotalAmount:    input.TotalAmount,
		Details:        input.Details,
		WeekDays:       input.WeekDays,
		Shift:          input.Shift,
		StartTime:      input.StartTime,
		HoursPerDay:    input.HoursPerDay,
		TotalHours:     input.TotalHours,
	}
	if err := s.repo.CreateContract(ctx, contract); err != nil {
		return nil, err
	}
	return contract, nil
}

func (s *ContractService) GetByID(ctx context.Context, id string) (*domain.Contract, error) {
	return s.repo.GetContractByID(ctx, id)
}

func (s *ContractService) Update(ctx context.Context, input ports.UpdateContractInput) (*domain.Contract, error) {
	contract, err := s.repo.GetContractByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if input.Status != nil {
		contract.Status = *input.Status
	}
	if input.HourRate != nil {
		contract.HourRate = *input.HourRate
	}
	if input.TotalAmount != nil {
		contract.TotalAmount = *input.TotalAmount
	}
	if input.Details != nil {
		contract.Details = input.Details
	}
	if input.WeekDays != nil {
		contract.WeekDays = input.WeekDays
	}
	if input.Shift != nil {
		contract.Shift = *input.Shift
	}
	if input.StartTime != nil {
		contract.StartTime = *input.StartTime
	}
	if input.HoursPerDay != nil {
		contract.HoursPerDay = *input.HoursPerDay
	}
	if input.TotalHours != nil {
		contract.TotalHours = *input.TotalHours
	}
	if err := s.repo.UpdateContract(ctx, contract); err != nil {
		return nil, err
	}
	return contract, nil
}

func (s *ContractService) Delete(ctx context.Context, id string) error {
	return s.repo.DeleteContract(ctx, id)
}

func (s *ContractService) GetByUserID(ctx context.Context, userID string) ([]*domain.Contract, error) {
	return s.repo.GetContractsByUserID(ctx, userID)
}

func (s *ContractService) GetAll(ctx context.Context) ([]*domain.Contract, error) {
	return s.repo.GetAllContracts(ctx)
}

func (s *ContractService) GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.Contract, error) {
	return s.repo.GetAllContractsByProfessionalID(ctx, professionalID)
}

func (s *ContractService) GetByStatus(ctx context.Context, status string) ([]*domain.Contract, error) {
	return s.repo.GetAllContractsByStatus(ctx, status)
}

func (s *ContractService) GetByUserIDAndStatus(ctx context.Context, userID string, status string) ([]*domain.Contract, error) {
	return s.repo.GetAllContractsByUserIDAndStatus(ctx, userID, status)
}

func (s *ContractService) GetByProfessionalIDAndStatus(ctx context.Context, professionalID string, status string) ([]*domain.Contract, error) {
	return s.repo.GetAAllContractsByProfessionalIDAndStatus(ctx, professionalID, status)
}

func (s *ContractService) GetByProfessionalCategory(ctx context.Context, category string) ([]*domain.Contract, error) {
	return s.repo.GetAAllContractsByProfessionalCategory(ctx, category)
}
