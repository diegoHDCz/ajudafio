package availability

import (
	"context"
	"errors"
	"testing"

	"github.com/diegoHDCz/ajudafio/internal/availability/domain"
	"github.com/diegoHDCz/ajudafio/internal/shared"
)

type mockAvailRepo struct {
	getByProfessionalID func(context.Context, string) ([]*domain.Availability, error)
	create              func(context.Context, *domain.Availability) (*domain.Availability, error)
	update              func(context.Context, *domain.Availability) (*domain.Availability, error)
	delete              func(context.Context, string) error
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

func makeAvailability() *domain.Availability {
	shifts := []shared.Shift{shared.ShiftMorning}
	return &domain.Availability{
		ID:             "avail-1",
		ProfessionalID: "prof-1",
		DayOfWeek:      []shared.WeekDay{shared.Monday, shared.Wednesday},
		Shift:          &shifts,
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
			if len(a.DayOfWeek) != len(input.DayOfWeek) {
				t.Errorf("DayOfWeek len: got %d, want %d", len(a.DayOfWeek), len(input.DayOfWeek))
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
