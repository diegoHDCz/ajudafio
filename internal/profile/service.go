package profile

import (
	"context"
	"fmt"

	"github.com/diegoHDCz/ajudafio/internal/profile/domain"
	"github.com/diegoHDCz/ajudafio/internal/profile/ports"
)

type service struct {
	repo ports.ProfileRepository
}

func NewService(repo ports.ProfileRepository) ports.ProfileService {
	return &service{repo: repo}
}

func (s *service) CreateFamilyProfile(ctx context.Context, userID string) (*domain.FamilyProfile, error) {
	p, err := s.repo.CreateFamilyProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("profile.CreateFamilyProfile: %w", err)
	}
	return p, nil
}

func (s *service) GetFamilyProfileByUserID(ctx context.Context, userID string) (*domain.FamilyProfile, error) {
	p, err := s.repo.GetFamilyProfileByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("profile.GetFamilyProfileByUserID: %w", err)
	}
	return p, nil
}

func (s *service) CreateFinancialProfile(ctx context.Context, userID string) (*domain.FinancialProfile, error) {
	p, err := s.repo.CreateFinancialProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("profile.CreateFinancialProfile: %w", err)
	}
	return p, nil
}

func (s *service) GetFinancialProfileByUserID(ctx context.Context, userID string) (*domain.FinancialProfile, error) {
	p, err := s.repo.GetFinancialProfileByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("profile.GetFinancialProfileByUserID: %w", err)
	}
	return p, nil
}
