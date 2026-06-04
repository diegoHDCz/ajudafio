package domain

import "github.com/diegoHDCz/ajudafio/internal/shared"

type Availability struct {
	ID             string
	ProfessionalID string
	DayOfWeek      shared.WeekDay
	Shift          *shared.Shift
	StartHour      *string
	EndHour        *string
}
