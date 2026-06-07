package http

import (
	"encoding/json"
	"time"

	"github.com/diegoHDCz/ajudafio/internal/professional/domain"
	professionalPorts "github.com/diegoHDCz/ajudafio/internal/professional/ports"
)

type createProfessionalRequest struct {
	UserID            string          `json:"user_id"`
	LicenseNumber     string          `json:"license_number"`
	Category          domain.Category `json:"category"`
	YearsOfExperience int             `json:"years_of_experience"`
	Resume            *string         `json:"resume"`
	Metadata          json.RawMessage `json:"metadata" swaggertype:"object"`
}

type updateProfessionalRequest struct {
	LicenseNumber     *string          `json:"license_number"`
	Category          *domain.Category `json:"category"`
	YearsOfExperience *int             `json:"years_of_experience"`
	Verified          *bool            `json:"verified"`
	Resume            *string          `json:"resume"`
	Metadata          json.RawMessage  `json:"metadata" swaggertype:"object"`
}

type professionalResponse struct {
	ID                string          `json:"id"`
	UserID            string          `json:"user_id"`
	LicenseNumber     string          `json:"license_number"`
	Category          domain.Category `json:"category"`
	YearsOfExperience int             `json:"years_of_experience"`
	Verified          bool            `json:"verified"`
	Resume            string          `json:"resume"`
	Metadata          json.RawMessage `json:"metadata,omitempty" swaggertype:"object"`
	CreatedAt         string          `json:"created_at"`
	UpdatedAt         string          `json:"updated_at"`
	Name              *string         `json:"name,omitempty"`
	AvatarURL         *string         `json:"avatar_url,omitempty"`
	Email             *string         `json:"email,omitempty"`
	Role              *string         `json:"role,omitempty"`
}

type professionalPageResponse struct {
	Items      []professionalResponse `json:"items"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}

func toPageResponse(page *professionalPorts.ProfessionalPage) professionalPageResponse {
	items := make([]professionalResponse, len(page.Items))
	for i, p := range page.Items {
		items[i] = toResponse(p)
	}
	return professionalPageResponse{
		Items:      items,
		Total:      page.Total,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalPages: page.TotalPages,
	}
}

func toResponse(p *domain.Professional) professionalResponse {
	var meta json.RawMessage
	if len(p.Metadata) > 0 {
		meta = json.RawMessage(p.Metadata)
	}
	return professionalResponse{
		ID:                p.ID,
		UserID:            p.UserID,
		LicenseNumber:     p.LicenseNumber,
		Category:          p.Category,
		YearsOfExperience: p.YearsOfExperience,
		Verified:          p.Verified,
		Resume:            p.Resume,
		Metadata:          meta,
		CreatedAt:         p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         p.UpdatedAt.Format(time.RFC3339),
		Name:              p.UserName,
		AvatarURL:         p.UserAvatarURL,
		Email:             p.UserEmail,
		Role:              p.UserRole,
	}
}
