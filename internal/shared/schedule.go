package shared

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
	ShiftFullDay   Shift = "FULL_DAY"
	ShiftCustom    Shift = "CUSTOM"
)

// ShiftHours maps a preset shift to its [startHour, endHour] in "HH:MM" format.
var ShiftHours = map[Shift][2]string{
	ShiftMorning:   {"09:00", "12:00"},
	ShiftAfternoon: {"13:00", "18:00"},
	ShiftNight:     {"19:00", "23:00"},
	ShiftFullDay:   {"08:00", "18:00"},
}
