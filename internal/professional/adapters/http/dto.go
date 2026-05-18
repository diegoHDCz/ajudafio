package http

import (
	"encoding/json"
	"time"

	"github.com/diegoHDCz/ajudafio/internal/professional/domain"
)

type createProfessionalRequest struct {
	UserID            string          `json:"user_id"`
	LicenseNumber     string          `json:"license_number"`
	Category          domain.Category `json:"category"`
	YearsOfExperience int             `json:"years_of_experience"`
	Resume            *string         `json:"resume"`
	Metadata          json.RawMessage `json:"metadata"`
}

type updateProfessionalRequest struct {
	LicenseNumber     *string          `json:"license_number"`
	Category          *domain.Category `json:"category"`
	YearsOfExperience *int             `json:"years_of_experience"`
	Verified          *bool            `json:"verified"`
	Resume            *string          `json:"resume"`
	Metadata          json.RawMessage  `json:"metadata"`
}

type professionalResponse struct {
	ID                string          `json:"id"`
	UserID            string          `json:"user_id"`
	LicenseNumber     string          `json:"license_number"`
	Category          domain.Category `json:"category"`
	YearsOfExperience int             `json:"years_of_experience"`
	Verified          bool            `json:"verified"`
	Resume            string          `json:"resume"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	CreatedAt         string          `json:"created_at"`
	UpdatedAt         string          `json:"updated_at"`
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
	}
}
