package availability

import (
	"context"
	"errors"
	"testing"

	"github.com/diegoHDCz/ajudafio/internal/availability/domain"
	"github.com/diegoHDCz/ajudafio/internal/shared"
)

type mockAvailRepo struct {
	getByID                func(context.Context, string) (*domain.Availability, error)
	getByProfessionalID    func(context.Context, string) ([]*domain.Availability, error)
	create                 func(context.Context, *domain.Availability) (*domain.Availability, error)
	update                 func(context.Context, *domain.Availability) (*domain.Availability, error)
	delete                 func(context.Context, string) error
	deleteByProfessionalID func(context.Context, string) error
}

func (m *mockAvailRepo) GetByID(ctx context.Context, id string) (*domain.Availability, error) {
	if m.getByID != nil {
		return m.getByID(ctx, id)
	}
	return nil, errors.New("not found")
}
func (m *mockAvailRepo) GetByProfessionalID(ctx context.Context, id string) ([]*domain.Availability, error) {
	return m.getByProfessionalID(ctx, id)
}
func (m *mockAvailRepo) Create(ctx context.Context, a *domain.Availability) (*domain.Availability, error) {
	return m.create(ctx, a)
}
func (m *mockAvailRepo) Update(ctx context.Context, a *domain.Availability) (*domain.Availability, error) {
	return m.update(ctx, a)
}
func (m *mockAvailRepo) Delete(ctx context.Context, id string) error {
	return m.delete(ctx, id)
}
func (m *mockAvailRepo) DeleteByProfessionalID(ctx context.Context, id string) error {
	if m.deleteByProfessionalID != nil {
		return m.deleteByProfessionalID(ctx, id)
	}
	return nil
}

func ptrShift(s shared.Shift) *shared.Shift { return &s }

func makeAvailability() *domain.Availability {
	return &domain.Availability{
		ID:             "avail-1",
		ProfessionalID: "prof-1",
		DayOfWeek:      shared.Monday,
		Shift:          ptrShift(shared.ShiftMorning),
	}
}

// --- GetByProfessionalID ---

func TestAvailabilityGetByProfessionalID_Success(t *testing.T) {
	list := []*domain.Availability{makeAvailability(), makeAvailability()}
	svc := NewAvailabilityService(&mockAvailRepo{
		getByProfessionalID: func(_ context.Context, id string) ([]*domain.Availability, error) {
			if id != "prof-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return list, nil
		},
	})

	got, err := svc.GetByProfessionalID(context.Background(), "prof-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 items, got %d", len(got))
	}
}

func TestAvailabilityGetByProfessionalID_Empty(t *testing.T) {
	svc := NewAvailabilityService(&mockAvailRepo{
		getByProfessionalID: func(_ context.Context, _ string) ([]*domain.Availability, error) {
			return []*domain.Availability{}, nil
		},
	})

	got, err := svc.GetByProfessionalID(context.Background(), "prof-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %d items", len(got))
	}
}

func TestAvailabilityGetByProfessionalID_Error(t *testing.T) {
	repoErr := errors.New("db error")
	svc := NewAvailabilityService(&mockAvailRepo{
		getByProfessionalID: func(_ context.Context, _ string) ([]*domain.Availability, error) {
			return nil, repoErr
		},
	})

	_, err := svc.GetByProfessionalID(context.Background(), "prof-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- Create ---

func TestAvailabilityCreate_Success(t *testing.T) {
	input := makeAvailability()
	want := makeAvailability()
	want.ID = "new-id"

	svc := NewAvailabilityService(&mockAvailRepo{
		create: func(_ context.Context, a *domain.Availability) (*domain.Availability, error) {
			if a.ProfessionalID != input.ProfessionalID {
				t.Errorf("ProfessionalID: got %s, want %s", a.ProfessionalID, input.ProfessionalID)
			}
			if a.DayOfWeek != input.DayOfWeek {
				t.Errorf("DayOfWeek: got %s, want %s", a.DayOfWeek, input.DayOfWeek)
			}
			if a.Shift == nil || input.Shift == nil || *a.Shift != *input.Shift {
				t.Errorf("Shift: got %v, want %v", a.Shift, input.Shift)
			}
			return want, nil
		},
	})

	got, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != want.ID {
		t.Errorf("ID: got %s, want %s", got.ID, want.ID)
	}
}

func TestAvailabilityCreate_Error(t *testing.T) {
	repoErr := errors.New("create failed")
	svc := NewAvailabilityService(&mockAvailRepo{
		create: func(_ context.Context, _ *domain.Availability) (*domain.Availability, error) {
			return nil, repoErr
		},
	})

	_, err := svc.Create(context.Background(), makeAvailability())
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- CreateBulk ---

func TestAvailabilityCreateBulk_Success(t *testing.T) {
	deleteCount := 0
	createCount := 0
	svc := NewAvailabilityService(&mockAvailRepo{
		deleteByProfessionalID: func(_ context.Context, _ string) error {
			deleteCount++
			return nil
		},
		create: func(_ context.Context, a *domain.Availability) (*domain.Availability, error) {
			createCount++
			return a, nil
		},
	})

	rules := []*domain.Availability{
		{ID: "1", ProfessionalID: "prof-1", DayOfWeek: shared.Monday, Shift: ptrShift(shared.ShiftMorning)},
		{ID: "2", ProfessionalID: "prof-1", DayOfWeek: shared.Tuesday, Shift: ptrShift(shared.ShiftAfternoon)},
	}
	got, err := svc.CreateBulk(context.Background(), "prof-1", rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleteCount != 1 {
		t.Errorf("deleteByProfessionalID called %d times, want 1", deleteCount)
	}
	if createCount != 2 {
		t.Errorf("create called %d times, want 2", createCount)
	}
	if len(got) != 2 {
		t.Errorf("returned %d items, want 2", len(got))
	}
}

func TestAvailabilityCreateBulk_ShiftResolution(t *testing.T) {
	svc := NewAvailabilityService(&mockAvailRepo{
		deleteByProfessionalID: func(_ context.Context, _ string) error { return nil },
		create: func(_ context.Context, a *domain.Availability) (*domain.Availability, error) {
			if a.StartHour == nil || *a.StartHour != "09:00" {
				t.Errorf("StartHour: got %v, want 09:00", a.StartHour)
			}
			if a.EndHour == nil || *a.EndHour != "12:00" {
				t.Errorf("EndHour: got %v, want 12:00", a.EndHour)
			}
			return a, nil
		},
	})

	rules := []*domain.Availability{
		{ID: "1", ProfessionalID: "prof-1", DayOfWeek: shared.Monday, Shift: ptrShift(shared.ShiftMorning)},
	}
	if _, err := svc.CreateBulk(context.Background(), "prof-1", rules); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAvailabilityCreateBulk_CustomRequiresHours(t *testing.T) {
	svc := NewAvailabilityService(&mockAvailRepo{})
	rules := []*domain.Availability{
		{ID: "1", ProfessionalID: "prof-1", DayOfWeek: shared.Monday},
	}
	if _, err := svc.CreateBulk(context.Background(), "prof-1", rules); err == nil {
		t.Error("expected error for custom rule without hours")
	}
}

func TestAvailabilityCreateBulk_OverlapRejected(t *testing.T) {
	start1, end1 := "10:00", "14:00"
	start2, end2 := "12:00", "16:00"
	svc := NewAvailabilityService(&mockAvailRepo{})
	rules := []*domain.Availability{
		{ID: "1", ProfessionalID: "prof-1", DayOfWeek: shared.Wednesday, StartHour: &start1, EndHour: &end1},
		{ID: "2", ProfessionalID: "prof-1", DayOfWeek: shared.Wednesday, StartHour: &start2, EndHour: &end2},
	}
	if _, err := svc.CreateBulk(context.Background(), "prof-1", rules); err == nil {
		t.Error("expected overlap error")
	}
}

func TestAvailabilityCreateBulk_NoOverlapDifferentDays(t *testing.T) {
	start, end := "10:00", "14:00"
	svc := NewAvailabilityService(&mockAvailRepo{
		deleteByProfessionalID: func(_ context.Context, _ string) error { return nil },
		create: func(_ context.Context, a *domain.Availability) (*domain.Availability, error) { return a, nil },
	})
	rules := []*domain.Availability{
		{ID: "1", ProfessionalID: "prof-1", DayOfWeek: shared.Monday, StartHour: &start, EndHour: &end},
		{ID: "2", ProfessionalID: "prof-1", DayOfWeek: shared.Tuesday, StartHour: &start, EndHour: &end},
	}
	if _, err := svc.CreateBulk(context.Background(), "prof-1", rules); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- Update ---

func TestAvailabilityUpdate_Success(t *testing.T) {
	input := makeAvailability()
	want := makeAvailability()

	svc := NewAvailabilityService(&mockAvailRepo{
		update: func(_ context.Context, a *domain.Availability) (*domain.Availability, error) {
			if a.ID != input.ID {
				t.Errorf("ID: got %s, want %s", a.ID, input.ID)
			}
			return want, nil
		},
	})

	got, err := svc.Update(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAvailabilityUpdate_Error(t *testing.T) {
	repoErr := errors.New("update failed")
	svc := NewAvailabilityService(&mockAvailRepo{
		update: func(_ context.Context, _ *domain.Availability) (*domain.Availability, error) {
			return nil, repoErr
		},
	})

	_, err := svc.Update(context.Background(), makeAvailability())
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- Delete ---

func TestAvailabilityDelete_Success(t *testing.T) {
	svc := NewAvailabilityService(&mockAvailRepo{
		delete: func(_ context.Context, id string) error {
			if id != "avail-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return nil
		},
	})

	if err := svc.Delete(context.Background(), "avail-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAvailabilityDelete_Error(t *testing.T) {
	repoErr := errors.New("delete failed")
	svc := NewAvailabilityService(&mockAvailRepo{
		delete: func(_ context.Context, _ string) error { return repoErr },
	})

	if err := svc.Delete(context.Background(), "avail-1"); !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}
