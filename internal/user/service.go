package user

import (
	"context"
	"fmt"

	"github.com/diegoHDCz/ajudafio/internal/user/domain"
	"github.com/diegoHDCz/ajudafio/internal/user/ports"
)

type service struct {
	repo ports.UserRepository
}

func NewService(repo ports.UserRepository) ports.UserService {
	return &service{repo: repo}
}

func (s *service) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user.GetByID: %w", err)
	}
	return user, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user.GetByEmail: %w", err)
	}
	return user, nil
}

func (s *service) Create(ctx context.Context, input ports.CreateUserInput) (*domain.User, error) {
	user := &domain.User{
		Email: input.Email,
		Name:  input.Name,  // Agora é string direta
		Phone: input.Phone, // Mapeado para o novo nome
		Role:  input.Role,
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("user.Create: %w", err)
	}
	return created, nil
}

func (s *service) Update(ctx context.Context, input ports.UpdateUserInput) (*domain.User, error) {
	existing, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("user.Update: %w", err)
	}

	// Atualizações parciais (Somente se o ponteiro não for nil)
	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Email != nil {
		existing.Email = *input.Email
	}
	if input.Phone != nil {
		existing.Phone = input.Phone
	}
	if input.Role != nil {
		existing.Role = *input.Role
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("user.Update: %w", err)
	}
	return updated, nil
}

func (s *service) Delete(ctx context.Context, id domain.UserID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("user.Delete: %w", err)
	}
	return nil
}
