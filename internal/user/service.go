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
		Email:                   input.Email,
		Name:                    input.Name,
		Telephone:               input.Telephone,
		TelephoneWhatsapp:       input.TelephoneWhatsapp,
		SecondTelephone:         input.SecondTelephone,
		SecondTelephoneWhatsapp: input.SecondTelephoneWhatsapp,
		Linkedin:                input.Linkedin,
		Instagram:               input.Instagram,
		Facebook:                input.Facebook,
		IdentificationNumber:    input.IdentificationNumber,
		IdentificationType:      input.IdentificationType,
		Role:                    input.Role,
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

	if input.Name != nil {
		existing.Name = input.Name
	}
	if input.Telephone != nil {
		existing.Telephone = input.Telephone
	}
	if input.TelephoneWhatsapp != nil {
		existing.TelephoneWhatsapp = *input.TelephoneWhatsapp
	}
	if input.SecondTelephone != nil {
		existing.SecondTelephone = input.SecondTelephone
	}
	if input.SecondTelephoneWhatsapp != nil {
		existing.SecondTelephoneWhatsapp = *input.SecondTelephoneWhatsapp
	}
	if input.Linkedin != nil {
		existing.Linkedin = input.Linkedin
	}
	if input.Instagram != nil {
		existing.Instagram = input.Instagram
	}
	if input.Facebook != nil {
		existing.Facebook = input.Facebook
	}
	if input.IdentificationNumber != nil {
		existing.IdentificationNumber = input.IdentificationNumber
	}
	if input.IdentificationType != nil {
		existing.IdentificationType = input.IdentificationType
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
