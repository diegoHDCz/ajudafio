package http

import (
	"encoding/json"
	"time"

	"github.com/diegoHDCz/ajudafio/internal/contract/domain"
)

type createContractRequest struct {
	ClientID       string          `json:"client_id"`
	ProfessionalID string          `json:"professional_id"`
	HourRate       int             `json:"hour_rate"`
	TotalAmount    int             `json:"total_amount"`
	Details        json.RawMessage `json:"details" swaggertype:"object"`
	WeekDays       []string        `json:"week_days"`
	Shift          string          `json:"shift"`
	StartTime      string          `json:"start_time"` // "15:04"
	HoursPerDay    int             `json:"hours_per_day"`
	TotalHours     int             `json:"total_hours"`
}

type updateContractRequest struct {
	Status      *string         `json:"status"`
	HourRate    *int            `json:"hour_rate"`
	TotalAmount *int            `json:"total_amount"`
	Details     json.RawMessage `json:"details" swaggertype:"object"`
	WeekDays    []string        `json:"week_days"`
	Shift       *string         `json:"shift"`
	StartTime   *string         `json:"start_time"` // "15:04"
	HoursPerDay *int            `json:"hours_per_day"`
	TotalHours  *int            `json:"total_hours"`
}

type contractResponse struct {
	ID             string          `json:"id"`
	ClientID       string          `json:"client_id"`
	ProfessionalID string          `json:"professional_id"`
	Status         string          `json:"status"`
	HourRate       int             `json:"hour_rate"`
	TotalAmount    int             `json:"total_amount"`
	Details        json.RawMessage `json:"details,omitempty" swaggertype:"object"`
	WeekDays       []string        `json:"week_days"`
	Shift          string          `json:"shift"`
	StartTime      string          `json:"start_time"`
	HoursPerDay    int             `json:"hours_per_day"`
	TotalHours     int             `json:"total_hours"`
	CreatedAt      string          `json:"created_at"`
}

func toResponse(c *domain.Contract) contractResponse {
	weekDays := make([]string, len(c.WeekDays))
	for i, d := range c.WeekDays {
		weekDays[i] = string(d)
	}
	return contractResponse{
		ID:             c.ID,
		ClientID:       c.ClientID,
		ProfessionalID: c.ProfessionalID,
		Status:         c.Status,
		HourRate:       c.HourRate,
		TotalAmount:    c.TotalAmount,
		Details:        json.RawMessage(c.Details),
		WeekDays:       weekDays,
		Shift:          string(c.Shift),
		StartTime:      c.StartTime.Format("15:04"),
		HoursPerDay:    c.HoursPerDay,
		TotalHours:     c.TotalHours,
		CreatedAt:      c.CreatedAt.Format(time.RFC3339),
	}
}
