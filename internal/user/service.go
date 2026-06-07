package user

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	storagePorts "github.com/diegoHDCz/ajudafio/internal/storage/ports"
	"github.com/diegoHDCz/ajudafio/internal/user/domain"
	"github.com/diegoHDCz/ajudafio/internal/user/ports"
	"github.com/google/uuid"
)

const avatarSignedURLExpiration = 15 * time.Minute

type service struct {
	repo    ports.UserRepository
	storage storagePorts.StorageProvider
}

func NewService(repo ports.UserRepository, storage storagePorts.StorageProvider) ports.UserService {
	return &service{repo: repo, storage: storage}
}

func (s *service) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user.GetByID: %w", err)
	}

	if user.AvatarURL != nil && s.storage != nil {
		if key := extractS3Key(*user.AvatarURL); key != "" {
			signedURL, err := s.storage.GetSignedURL(ctx, key, avatarSignedURLExpiration)
			if err != nil {
				return nil, fmt.Errorf("user.GetByID: %w", err)
			}
			user.AvatarURL = &signedURL
		}
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
		ID:    uuid.New().String(),
		Email: input.Email,
		Name:  input.Name,
		Phone: input.Phone,
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

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("user.Delete: %w", err)
	}
	return nil
}

func (s *service) UpdateUserRole(ctx context.Context, id string, role domain.Role) error {
	if err := s.repo.UpdateUserRole(ctx, id, role); err != nil {
		return fmt.Errorf("user.UpdateUserRole: %w", err)
	}
	return nil
}

func (s *service) UploadAvatar(ctx context.Context, userID string, fileData []byte, contentType string) (*domain.User, error) {
	if s.storage == nil {
		return nil, fmt.Errorf("user.UploadAvatar: storage not configured")
	}

	existing, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user.UploadAvatar: %w", err)
	}

	if existing.AvatarURL != nil {
		if key := extractS3Key(*existing.AvatarURL); key != "" {
			_ = s.storage.Delete(ctx, key)
		}
	}

	ext := extensionFromContentType(contentType)
	key := fmt.Sprintf("avatars/%s-%d%s", userID, time.Now().UnixMilli(), ext)

	newURL, err := s.storage.Upload(ctx, key, fileData, contentType)
	if err != nil {
		return nil, fmt.Errorf("user.UploadAvatar: %w", err)
	}

	updated, err := s.repo.UpdateAvatar(ctx, userID, &newURL)
	if err != nil {
		return nil, fmt.Errorf("user.UploadAvatar: %w", err)
	}
	return updated, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func extractS3Key(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(u.Path, "/")
}

func extensionFromContentType(ct string) string {
	switch ct {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		return ".bin"
	}
}
