package domain

import (
	"time"
)

type WeekDay string

const (
	Monday    WeekDay = "MONDAY"
	Tuesday   WeekDay = "TUESDAY"
	Wednesday WeekDay = "WEDNESDAY"
	Thursday  WeekDay = "THURSDAY"
	Friday    WeekDay = "FRIDAY"
	Saturday  WeekDay = "SATURDAY"
	Sunday    WeekDay = "SUNDAY"
)

type Shift string

const (
	ShiftMorning   Shift = "MORNING"
	ShiftAfternoon Shift = "AFTERNOON"
	ShiftNight     Shift = "NIGHT"
	ShiftCustom    Shift = "CUSTOM"
)

type Contract struct {
	ID             string
	ClientID       string
	ProfessionalID string
	Status         string
	HourRate       int
	TotalAmount    int
	Details        []byte
	WeekDays       []WeekDay
	Shift          Shift
	StartTime      time.Time
	HoursPerDay    int
	TotalHours     int
	CreatedAt      time.Time
}
