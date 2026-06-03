package domain

import "time"

type Appointment struct {
	ID             string
	ContractID     string
	ClientID       string
	ProfessionalID string
	Date           time.Time
	StartTime      time.Time
	EndTime        time.Time
	Status         string
	ZipCode        string
	AddressLine    string
	Number         string
	Complement     *string
	District       string
	City           string
	State          string
	Reference      *string
	Version        int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
