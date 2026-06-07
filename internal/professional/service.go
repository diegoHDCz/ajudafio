package professional

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/diegoHDCz/ajudafio/internal/professional/domain"
	"github.com/diegoHDCz/ajudafio/internal/professional/ports"
	storagePorts "github.com/diegoHDCz/ajudafio/internal/storage/ports"
	"github.com/google/uuid"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100

	avatarSignedURLExpiration = 15 * time.Minute
)

type ProfessionalService struct {
	ps      ports.ProfessionalRepository
	storage storagePorts.StorageProvider
}

func NewProfessionalService(ps ports.ProfessionalRepository, storage storagePorts.StorageProvider) *ProfessionalService {
	return &ProfessionalService{ps: ps, storage: storage}
}

func (s *ProfessionalService) GetByID(ctx context.Context, id string) (*domain.Professional, error) {
	professional, err := s.ps.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return professional, nil
}

func (s *ProfessionalService) GetByUserID(ctx context.Context, userID string) (*domain.Professional, error) {
	professional, err := s.ps.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return professional, nil
}

func (s *ProfessionalService) Create(ctx context.Context, input ports.CreateProfessionalInput) (*domain.Professional, error) {
	resume := ""
	if input.Resume != nil {
		resume = *input.Resume
	}
	professional, err := domain.NewProfessional(
		uuid.New().String(),
		input.UserID,
		input.LicenseNumber,
		input.Category,
		input.YearsOfExperience,
		resume,
		input.Metadata,
	)
	if err != nil {
		return nil, err
	}
	return s.ps.Create(ctx, professional)
}

func (s *ProfessionalService) Update(ctx context.Context, input ports.UpdateProfessionalInput) (*domain.Professional, error) {
	professional, err := s.ps.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if err := professional.ApplyUpdate(
		input.LicenseNumber,
		input.Category,
		input.YearsOfExperience,
		input.Verified,
		input.Resume,
		input.Metadata,
	); err != nil {
		return nil, err
	}

	return s.ps.Update(ctx, professional)
}

func (s *ProfessionalService) Delete(ctx context.Context, id string) error {
	return s.ps.Delete(ctx, id)
}

func (s *ProfessionalService) FindWithFilters(ctx context.Context, filters ports.ProfessionalFilters) (*ports.ProfessionalPage, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = defaultPageSize
	}
	if filters.PageSize > maxPageSize {
		filters.PageSize = maxPageSize
	}

	items, total, err := s.ps.FindWithFilters(ctx, filters)
	if err != nil {
		return nil, err
	}

	if s.storage != nil {
		for _, p := range items {
			if p.UserAvatarURL == nil {
				continue
			}
			key := extractS3Key(*p.UserAvatarURL)
			if key == "" {
				continue
			}
			signedURL, err := s.storage.GetSignedURL(ctx, key, avatarSignedURLExpiration)
			if err != nil {
				return nil, fmt.Errorf("professional.FindWithFilters: %w", err)
			}
			p.UserAvatarURL = &signedURL
		}
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))
	return &ports.ProfessionalPage{
		Items:      items,
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func extractS3Key(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(u.Path, "/")
}
