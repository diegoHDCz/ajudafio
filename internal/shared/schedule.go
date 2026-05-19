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
	ShiftCustom    Shift = "CUSTOM"
)
