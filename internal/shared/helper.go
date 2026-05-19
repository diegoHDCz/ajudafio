package shared

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ParseUUID(s string) (pgtype.UUID, error) {
	uid, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: uid, Valid: true}, nil
}

func ParseDayOfWeek(days []WeekDay) ([]string, error) {

	if len(days) == 0 {
		return nil, errors.New("a lista de dias não pode estar vazia")
	}

	// Converte o slice de tipos customizados (DayOfWeek) para um slice de string padrão
	var stringDays []string
	for _, day := range days {
		stringDays = append(stringDays, string(day))
	}

	return stringDays, nil
}

func SliceToDayOfWeek(days []string) []WeekDay {
	weekDays := make([]WeekDay, len(days))
	for i, day := range days {
		weekDays[i] = WeekDay(day)
	}
	return weekDays
}

func FormatDayOfWeek(days string) ([]WeekDay, error) {

	days = strings.TrimSpace(days)

	if days == "" {
		return nil, errors.New("a string de dias não pode estar vazia")
	}

	parts := strings.Split(days, ",")
	var weekDays []WeekDay
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if trimmedPart != "" {
			weekDays = append(weekDays, WeekDay(trimmedPart))
		}
	}
	if len(weekDays) == 0 {
		return nil, errors.New("nenhum dia válido encontrado na string")
	}
	return weekDays, nil
}
