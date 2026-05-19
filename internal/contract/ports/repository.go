package ports

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/contract/domain"
)

type ContractRepository interface {
	CreateContract(context.Context, *domain.Contract) error
	GetContractByID(context.Context, string) (*domain.Contract, error)
	UpdateContract(context.Context, *domain.Contract) error
	DeleteContract(context.Context, string) error
	GetContractsByUserID(context.Context, string) ([]*domain.Contract, error)
	GetAllContracts(context.Context) ([]*domain.Contract, error)
	GetAllContractsByProfessionalID(context.Context, string) ([]*domain.Contract, error)
	GetAllContractsByStatus(context.Context, string) ([]*domain.Contract, error)
	GetAllContractsByUserIDAndStatus(context.Context, string, string) ([]*domain.Contract, error)
	GetAAllContractsByProfessionalIDAndStatus(context.Context, string, string) ([]*domain.Contract, error)
	GetAAllContractsByProfessionalCategory(context.Context, string) ([]*domain.Contract, error)
}
