package http

import (
	"encoding/json"
	"net/http"

	"github.com/diegoHDCz/ajudafio/internal/appointment/domain"
)

type createAppointmentRequest struct {
	ContractID     string  `json:"contract_id"`
	ClientID       string  `json:"client_id"`
	ProfessionalID string  `json:"professional_id"`
	Date           string  `json:"date"`       // "2006-01-02"
	StartTime      string  `json:"start_time"` // "15:04"
	EndTime        string  `json:"end_time"`   // "15:04"
	ZipCode        string  `json:"zip_code"`
	AddressLine    string  `json:"address_line"`
	Number         string  `json:"number"`
	Complement     *string `json:"complement"`
	District       string  `json:"district"`
	City           string  `json:"city"`
	State          string  `json:"state"`
	Reference      *string `json:"reference"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

type appointmentResponse struct {
	ID             string  `json:"id"`
	ContractID     string  `json:"contract_id"`
	ClientID       string  `json:"client_id"`
	ProfessionalID string  `json:"professional_id"`
	Date           string  `json:"date"`
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	Status         string  `json:"status"`
	ZipCode        string  `json:"zip_code"`
	AddressLine    string  `json:"address_line"`
	Number         string  `json:"number"`
	Complement     *string `json:"complement"`
	District       string  `json:"district"`
	City           string  `json:"city"`
	State          string  `json:"state"`
	Reference      *string `json:"reference"`
	Version        int     `json:"version"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

func toResponse(a *domain.Appointment) appointmentResponse {
	return appointmentResponse{
		ID:             a.ID,
		ContractID:     a.ContractID,
		ClientID:       a.ClientID,
		ProfessionalID: a.ProfessionalID,
		Date:           a.Date.Format("2006-01-02"),
		StartTime:      a.StartTime.Format("15:04"),
		EndTime:        a.EndTime.Format("15:04"),
		Status:         a.Status,
		ZipCode:        a.ZipCode,
		AddressLine:    a.AddressLine,
		Number:         a.Number,
		Complement:     a.Complement,
		District:       a.District,
		City:           a.City,
		State:          a.State,
		Reference:      a.Reference,
		Version:        a.Version,
		CreatedAt:      a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      a.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
