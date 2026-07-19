package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	"github.com/diegoHDCz/ajudafio/internal/auth/ports"
	userdomain "github.com/diegoHDCz/ajudafio/internal/user/domain"
	userports "github.com/diegoHDCz/ajudafio/internal/user/ports"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrEmailAlreadyInUse   = errors.New("email already in use")
	ErrRefreshTokenInvalid = errors.New("refresh token is no longer valid")
)

type service struct {
	repo    ports.AuthRepository
	userSvc userports.UserService
	secret  []byte
}

func NewService(repo ports.AuthRepository, userSvc userports.UserService, secret []byte) ports.AuthService {
	return &service{repo: repo, userSvc: userSvc, secret: secret}
}

func (s *service) Register(ctx context.Context, input ports.RegisterInput) (*authdomain.TokenPair, error) {
	if _, err := s.userSvc.GetByEmail(ctx, input.Email); err == nil {
		return nil, ErrEmailAlreadyInUse
	}

	hash, err := HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("auth.Register: %w", err)
	}

	// Role is always CLIENT on self-registration; only an admin can promote a user afterwards.
	user, err := s.userSvc.Create(ctx, userports.CreateUserInput{
		ID:    uuid.New().String(),
		Email: input.Email,
		Name:  input.Name,
		Phone: input.Phone,
		Role:  userdomain.RoleClient,
	})
	if err != nil {
		return nil, fmt.Errorf("auth.Register: %w", err)
	}

	if err := s.repo.SetPasswordHash(ctx, user.ID, hash); err != nil {
		return nil, fmt.Errorf("auth.Register: %w", err)
	}

	return s.issueTokenPair(ctx, user.ID, user.Email, user.Name, string(user.Role))
}

func (s *service) Login(ctx context.Context, email, password string) (*authdomain.TokenPair, error) {
	userAuth, err := s.repo.GetUserAuthByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if userAuth.PasswordHash == "" || !VerifyPassword(userAuth.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}
	return s.issueTokenPair(ctx, userAuth.ID, userAuth.Email, userAuth.Name, userAuth.Role)
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (*authdomain.TokenPair, error) {
	claims, err := ValidateRefreshToken(s.secret, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("auth.Refresh: %w", err)
	}

	tokenHash := HashToken(refreshToken)
	stored, err := s.repo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("auth.Refresh: %w", err)
	}
	if stored.RevokedAt != nil || time.Now().After(stored.ExpiresAt) {
		return nil, ErrRefreshTokenInvalid
	}

	user, err := s.userSvc.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("auth.Refresh: %w", err)
	}

	// Rotation: the token just used to refresh is revoked immediately, a new one is issued.
	if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return nil, fmt.Errorf("auth.Refresh: %w", err)
	}

	return s.issueTokenPair(ctx, user.ID, user.Email, user.Name, string(user.Role))
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	if err := s.repo.RevokeRefreshToken(ctx, HashToken(refreshToken)); err != nil {
		return fmt.Errorf("auth.Logout: %w", err)
	}
	return nil
}

func (s *service) issueTokenPair(ctx context.Context, userID, email, name, role string) (*authdomain.TokenPair, error) {
	access, err := GenerateAccessToken(s.secret, userID, email, name, role)
	if err != nil {
		return nil, fmt.Errorf("issueTokenPair: %w", err)
	}

	refresh, expiresAt, err := GenerateRefreshToken(s.secret, userID)
	if err != nil {
		return nil, fmt.Errorf("issueTokenPair: %w", err)
	}

	if _, err := s.repo.CreateRefreshToken(ctx, &authdomain.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		TokenHash: HashToken(refresh),
		ExpiresAt: expiresAt,
	}); err != nil {
		return nil, fmt.Errorf("issueTokenPair: %w", err)
	}

	return &authdomain.TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int(AccessTokenTTL.Seconds()),
	}, nil
}
