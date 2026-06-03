package http

import (
	"time"

	"github.com/diegoHDCz/ajudafio/internal/address/domain"
)

type createAddressRequest struct {
	UserID      string  `json:"user_id"`
	ZipCode     string  `json:"zip_code"`
	AddressLine string  `json:"address_line"`
	Number      string  `json:"number"`
	Complement  *string `json:"complement"`
	District    string  `json:"district"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Reference   *string `json:"reference"`
}

type updateAddressRequest struct {
	ZipCode     *string `json:"zip_code"`
	AddressLine *string `json:"address_line"`
	Number      *string `json:"number"`
	Complement  *string `json:"complement"`
	District    *string `json:"district"`
	City        *string `json:"city"`
	State       *string `json:"state"`
	Reference   *string `json:"reference"`
}

type addressResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	ZipCode     string  `json:"zip_code"`
	AddressLine string  `json:"address_line"`
	Number      string  `json:"number"`
	Complement  *string `json:"complement,omitempty"`
	District    string  `json:"district"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Reference   *string `json:"reference,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func toResponse(a *domain.Address) addressResponse {
	return addressResponse{
		ID:          a.ID,
		UserID:      a.UserID,
		ZipCode:     a.ZipCode,
		AddressLine: a.AddressLine,
		Number:      a.Number,
		Complement:  a.Complement,
		District:    a.District,
		City:        a.City,
		State:       a.State,
		Reference:   a.Reference,
		CreatedAt:   a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   a.UpdatedAt.Format(time.RFC3339),
	}
}
