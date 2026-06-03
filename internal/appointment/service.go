package appointment

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/diegoHDCz/ajudafio/internal/appointment/domain"
	"github.com/diegoHDCz/ajudafio/internal/appointment/ports"
	availabilityports "github.com/diegoHDCz/ajudafio/internal/availability/ports"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	"github.com/google/uuid"
)

var _ ports.AppointmentService = (*AppointmentService)(nil)

type AppointmentService struct {
	repo         ports.AppointmentRepository
	availability availabilityports.AvailabilityRepository
}

func NewAppointmentService(repo ports.AppointmentRepository, availability availabilityports.AvailabilityRepository) *AppointmentService {
	return &AppointmentService{repo: repo, availability: availability}
}

func (s *AppointmentService) Create(ctx context.Context, input ports.CreateAppointmentInput) (*domain.Appointment, error) {
	if err := s.validateAvailability(ctx, input); err != nil {
		return nil, err
	}

	overlap, err := s.repo.HasOverlap(ctx, input.ProfessionalID, input.Date, input.StartTime, input.EndTime)
	if err != nil {
		return nil, fmt.Errorf("appointment.Create: %w", err)
	}
	if overlap {
		return nil, errors.New("profissional já possui agendamento neste horário")
	}

	a := &domain.Appointment{
		ID:             uuid.New().String(),
		ContractID:     input.ContractID,
		ClientID:       input.ClientID,
		ProfessionalID: input.ProfessionalID,
		Date:           input.Date,
		StartTime:      input.StartTime,
		EndTime:        input.EndTime,
		ZipCode:        input.ZipCode,
		AddressLine:    input.AddressLine,
		Number:         input.Number,
		Complement:     input.Complement,
		District:       input.District,
		City:           input.City,
		State:          input.State,
		Reference:      input.Reference,
	}
	if err := s.repo.CreateAppointment(ctx, a); err != nil {
		return nil, fmt.Errorf("appointment.Create: %w", err)
	}
	return a, nil
}

func (s *AppointmentService) GetByID(ctx context.Context, id string) (*domain.Appointment, error) {
	return s.repo.GetAppointmentByID(ctx, id)
}

func (s *AppointmentService) GetByContractID(ctx context.Context, contractID string) ([]*domain.Appointment, error) {
	return s.repo.GetAppointmentsByContractID(ctx, contractID)
}

func (s *AppointmentService) GetByClientID(ctx context.Context, clientID string) ([]*domain.Appointment, error) {
	return s.repo.GetAppointmentsByClientID(ctx, clientID)
}

func (s *AppointmentService) GetByProfessionalID(ctx context.Context, professionalID string) ([]*domain.Appointment, error) {
	return s.repo.GetAppointmentsByProfessionalID(ctx, professionalID)
}

func (s *AppointmentService) UpdateStatus(ctx context.Context, id, status string) (*domain.Appointment, error) {
	current, err := s.repo.GetAppointmentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.repo.UpdateAppointmentStatus(ctx, id, status, current.Version)
}

func (s *AppointmentService) Delete(ctx context.Context, id string) error {
	return s.repo.DeleteAppointment(ctx, id)
}

func (s *AppointmentService) validateAvailability(ctx context.Context, input ports.CreateAppointmentInput) error {
	dayOfWeek := weekdayToString(input.Date.Weekday())

	availabilities, err := s.availability.GetByProfessionalID(ctx, input.ProfessionalID)
	if err != nil {
		return fmt.Errorf("appointment.validateAvailability: %w", err)
	}

	for _, av := range availabilities {
		if string(av.DayOfWeek) != dayOfWeek {
			continue
		}
		if coversTime(av.Shift, av.StartHour, av.EndHour, input.StartTime, input.EndTime) {
			return nil
		}
	}

	return fmt.Errorf("profissional não possui disponibilidade na %s para o horário %s–%s",
		dayOfWeek,
		input.StartTime.Format("15:04"),
		input.EndTime.Format("15:04"),
	)
}

func coversTime(shift shared.Shift, startHour, endHour *string, start, end time.Time) bool {
	switch shift {
	case shared.ShiftFullDay:
		return true
	case shared.ShiftMorning:
		return !start.Before(hm(6, 0)) && !end.After(hm(12, 0))
	case shared.ShiftAfternoon:
		return !start.Before(hm(12, 0)) && !end.After(hm(18, 0))
	case shared.ShiftNight:
		return !start.Before(hm(18, 0))
	case shared.ShiftCustom:
		if startHour == nil || endHour == nil {
			return false
		}
		avStart, err1 := time.Parse("15:04", *startHour)
		avEnd, err2 := time.Parse("15:04", *endHour)
		if err1 != nil || err2 != nil {
			return false
		}
		return !start.Before(avStart) && !end.After(avEnd)
	}
	return false
}

func hm(hour, minute int) time.Time {
	return time.Date(0, 1, 1, hour, minute, 0, 0, time.UTC)
}

func weekdayToString(w time.Weekday) string {
	return strings.ToUpper(w.String())
}
