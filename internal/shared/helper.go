package shared

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// NormalizeEmail trims whitespace, lowercases and validates the address format,
// so that the users.email UNIQUE constraint (case-sensitive in Postgres) can't
// be bypassed by casing/whitespace variants of the same address.
func NormalizeEmail(raw string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(raw))
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return "", errors.New("invalid email")
	}

	domain := email[strings.LastIndex(email, "@")+1:]
	if !strings.Contains(domain, ".") {
		return "", errors.New("invalid email")
	}

	return email, nil
}

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
