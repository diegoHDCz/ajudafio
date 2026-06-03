package http

import "time"

type createBookingRequestRequest struct {
	ClientID        string                  `json:"client_id"`
	ProfessionalID  string                  `json:"professional_id"`
	AddressID       string                  `json:"address_id"`
	ProposedValue   float64                 `json:"proposed_value"`
	ScheduleDetails scheduleDetailsRequest  `json:"schedule_details"`
}

type scheduleDetailsRequest struct {
	Recurrence []scheduleEntryRequest `json:"recurrence"`
	StartDate  string                 `json:"start_date"`
	EndDate    string                 `json:"end_date"`
}

type scheduleEntryRequest struct {
	DayOfWeek string `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type updateStatusRequest struct {
	Status          string  `json:"status"`
	RejectionReason *string `json:"rejection_reason,omitempty"`
}

type bookingRequestResponse struct {
	ID              string                   `json:"id"`
	ClientID        string                   `json:"client_id"`
	ProfessionalID  string                   `json:"professional_id"`
	AddressID       string                   `json:"address_id"`
	ProposedValue   float64                  `json:"proposed_value"`
	ScheduleDetails scheduleDetailsResponse  `json:"schedule_details"`
	Status          string                   `json:"status"`
	RejectionReason *string                  `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time                `json:"created_at"`
	RespondedAt     *time.Time               `json:"responded_at,omitempty"`
}

type scheduleDetailsResponse struct {
	Recurrence []scheduleEntryResponse `json:"recurrence"`
	StartDate  string                  `json:"start_date"`
	EndDate    string                  `json:"end_date"`
}

type scheduleEntryResponse struct {
	DayOfWeek string `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
