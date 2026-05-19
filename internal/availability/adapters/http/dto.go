package http

import "github.com/diegoHDCz/ajudafio/internal/shared"

type createAvailabilityRequest struct {
	ProfessionalID string           `json:"professional_id"`
	DayOfWeek      []shared.WeekDay `json:"day_of_week"`
	Shift          *[]shared.Shift  `json:"shift,omitempty"`
	StartHour      *string          `json:"start_hour,omitempty"`
	EndHour        *string          `json:"end_hour,omitempty"`
}

type updateAvailabilityRequest struct {
	DayOfWeek []shared.WeekDay `json:"day_of_week"`
	Shift     *[]shared.Shift  `json:"shift,omitempty"`
	StartHour *string          `json:"start_hour,omitempty"`
	EndHour   *string          `json:"end_hour,omitempty"`
}

type availabilityResponse struct {
	ID             string           `json:"id"`
	ProfessionalID string           `json:"professional_id"`
	DayOfWeek      []shared.WeekDay `json:"day_of_week"`
	Shift          *[]shared.Shift  `json:"shift,omitempty"`
	StartHour      *string          `json:"start_hour,omitempty"`
	EndHour        *string          `json:"end_hour,omitempty"`
}
