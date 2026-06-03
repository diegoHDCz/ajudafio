package domain

import "time"

type BookingRequestStatus string

const (
	StatusPending   BookingRequestStatus = "PENDING"
	StatusAccepted  BookingRequestStatus = "ACCEPTED"
	StatusRejected  BookingRequestStatus = "REJECTED"
	StatusExpired   BookingRequestStatus = "EXPIRED"
	StatusCancelled BookingRequestStatus = "CANCELLED"
)

type ScheduleEntry struct {
	DayOfWeek string `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type ScheduleDetails struct {
	Recurrence []ScheduleEntry `json:"recurrence"`
	StartDate  string          `json:"start_date"`
	EndDate    string          `json:"end_date"`
}

type BookingRequest struct {
	ID              string
	ClientID        string
	ProfessionalID  string
	AddressID       string
	ProposedValue   float64
	ScheduleDetails ScheduleDetails
	Status          BookingRequestStatus
	RejectionReason *string
	CreatedAt       time.Time
	RespondedAt     *time.Time
}
