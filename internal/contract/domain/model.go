package domain

import (
	"time"

	"github.com/diegoHDCz/ajudafio/internal/shared"
)

type Contract struct {
	ID             string
	ClientID       string
	ProfessionalID string
	Status         string
	HourRate       int
	TotalAmount    int
	Details        []byte
	WeekDays       []shared.WeekDay
	Shift          shared.Shift
	StartTime      time.Time
	HoursPerDay    int
	TotalHours     int
	CreatedAt      time.Time
}
