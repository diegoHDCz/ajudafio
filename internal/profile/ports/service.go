package ports

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/profile/domain"
)

type ProfileService interface {
	CreateFamilyProfile(ctx context.Context, userID string) (*domain.FamilyProfile, error)
	GetFamilyProfileByUserID(ctx context.Context, userID string) (*domain.FamilyProfile, error)
	CreateFinancialProfile(ctx context.Context, userID string) (*domain.FinancialProfile, error)
	GetFinancialProfileByUserID(ctx context.Context, userID string) (*domain.FinancialProfile, error)
}
