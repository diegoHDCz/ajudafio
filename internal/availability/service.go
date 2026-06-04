package availability

import (
	"context"
	"fmt"

	"github.com/diegoHDCz/ajudafio/internal/availability/domain"
	"github.com/diegoHDCz/ajudafio/internal/availability/ports"
	"github.com/diegoHDCz/ajudafio/internal/shared"
)

type AvailabilityService struct {
	ar ports.AvailabilityRepository
}

func NewAvailabilityService(ar ports.AvailabilityRepository) *AvailabilityService {
	return &AvailabilityService{ar: ar}
}

func (s *AvailabilityService) GetByID(ctx context.Context, id string) (*domain.Availability, error) {
	return s.ar.GetByID(ctx, id)
}

func (s *AvailabilityService) GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.Availability, error) {
	return s.ar.GetByProfessionalID(ctx, professionalID)
}

func (s *AvailabilityService) Create(ctx context.Context, availability *domain.Availability) (*domain.Availability, error) {
	return s.ar.Create(ctx, availability)
}

func (s *AvailabilityService) Update(ctx context.Context, availability *domain.Availability) (*domain.Availability, error) {
	return s.ar.Update(ctx, availability)
}

func (s *AvailabilityService) Delete(ctx context.Context, id string) error {
	return s.ar.Delete(ctx, id)
}

// CreateBulk replaces all existing availabilities for the professional with the provided rules.
// Preset shifts are mapped to canonical start/end hours. Overlapping intervals are rejected.
func (s *AvailabilityService) CreateBulk(ctx context.Context, professionalID string, rules []*domain.Availability) ([]*domain.Availability, error) {
	if err := resolveAndValidate(rules); err != nil {
		return nil, err
	}

	if err := s.ar.DeleteByProfessionalID(ctx, professionalID); err != nil {
		return nil, fmt.Errorf("availability.CreateBulk: %w", err)
	}

	created := make([]*domain.Availability, 0, len(rules))
	for _, rule := range rules {
		a, err := s.ar.Create(ctx, rule)
		if err != nil {
			return nil, fmt.Errorf("availability.CreateBulk: %w", err)
		}
		created = append(created, a)
	}
	return created, nil
}

// resolveAndValidate maps preset shifts to hours and checks for overlapping intervals per day.
func resolveAndValidate(rules []*domain.Availability) error {
	for _, rule := range rules {
		if rule.Shift != nil {
			hours, ok := shared.ShiftHours[*rule.Shift]
			if !ok {
				return fmt.Errorf("turno inválido: %s", *rule.Shift)
			}
			rule.StartHour = &hours[0]
			rule.EndHour = &hours[1]
		} else {
			if rule.StartHour == nil || rule.EndHour == nil {
				return fmt.Errorf("start_hour e end_hour são obrigatórios quando shift não é informado")
			}
		}
	}
	return checkOverlap(rules)
}

// checkOverlap returns an error if any two rules on the same day have overlapping intervals.
// Times are "HH:MM" strings which compare correctly lexicographically.
func checkOverlap(rules []*domain.Availability) error {
	byDay := make(map[shared.WeekDay][]*domain.Availability)
	for _, r := range rules {
		byDay[r.DayOfWeek] = append(byDay[r.DayOfWeek], r)
	}
	for day, group := range byDay {
		for i := 0; i < len(group); i++ {
			for j := i + 1; j < len(group); j++ {
				a, b := group[i], group[j]
				if *a.StartHour < *b.EndHour && *b.StartHour < *a.EndHour {
					return fmt.Errorf("horários sobrepostos em %s: %s–%s e %s–%s",
						day, *a.StartHour, *a.EndHour, *b.StartHour, *b.EndHour)
				}
			}
		}
	}
	return nil
}
