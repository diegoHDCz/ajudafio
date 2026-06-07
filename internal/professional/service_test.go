package professional

import (
	"context"
	"errors"
	"testing"

	"github.com/diegoHDCz/ajudafio/internal/professional/domain"
	"github.com/diegoHDCz/ajudafio/internal/professional/ports"
)

// mockProfessionalRepo implements ports.ProfessionalRepository for testing.
type mockProfessionalRepo struct {
	getByID         func(context.Context, string) (*domain.Professional, error)
	getByUserID     func(context.Context, string) (*domain.Professional, error)
	create          func(context.Context, *domain.Professional) (*domain.Professional, error)
	update          func(context.Context, *domain.Professional) (*domain.Professional, error)
	delete          func(context.Context, string) error
	findWithFilters func(context.Context, ports.ProfessionalFilters) ([]*domain.Professional, int64, error)
}

func (m *mockProfessionalRepo) GetByID(ctx context.Context, id string) (*domain.Professional, error) {
	return m.getByID(ctx, id)
}
func (m *mockProfessionalRepo) GetByUserID(ctx context.Context, userID string) (*domain.Professional, error) {
	return m.getByUserID(ctx, userID)
}
func (m *mockProfessionalRepo) Create(ctx context.Context, p *domain.Professional) (*domain.Professional, error) {
	return m.create(ctx, p)
}
func (m *mockProfessionalRepo) Update(ctx context.Context, p *domain.Professional) (*domain.Professional, error) {
	return m.update(ctx, p)
}
func (m *mockProfessionalRepo) Delete(ctx context.Context, id string) error {
	return m.delete(ctx, id)
}
func (m *mockProfessionalRepo) FindWithFilters(ctx context.Context, f ports.ProfessionalFilters) ([]*domain.Professional, int64, error) {
	return m.findWithFilters(ctx, f)
}

func ptr[T any](v T) *T { return &v }

func makeProfessional() *domain.Professional {
	return &domain.Professional{
		ID:                "prof-1",
		UserID:            "user-1",
		LicenseNumber:     "LIC-001",
		Category:          domain.Nurse,
		YearsOfExperience: 5,
		Verified:          false,
	}
}

// --- GetByID ---

func TestProfessionalGetByID_Success(t *testing.T) {
	want := makeProfessional()
	svc := NewProfessionalService(&mockProfessionalRepo{
		getByID: func(_ context.Context, id string) (*domain.Professional, error) {
			if id != want.ID {
				t.Fatalf("unexpected id: %s", id)
			}
			return want, nil
		},
	}, nil)

	got, err := svc.GetByID(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestProfessionalGetByID_RepoError(t *testing.T) {
	repoErr := errors.New("not found")
	svc := NewProfessionalService(&mockProfessionalRepo{
		getByID: func(_ context.Context, _ string) (*domain.Professional, error) { return nil, repoErr },
	}, nil)

	_, err := svc.GetByID(context.Background(), "prof-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- GetByUserID ---

func TestProfessionalGetByUserID_Success(t *testing.T) {
	want := makeProfessional()
	svc := NewProfessionalService(&mockProfessionalRepo{
		getByUserID: func(_ context.Context, userID string) (*domain.Professional, error) {
			if userID != want.UserID {
				t.Fatalf("unexpected userID: %s", userID)
			}
			return want, nil
		},
	}, nil)

	got, err := svc.GetByUserID(context.Background(), want.UserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestProfessionalGetByUserID_RepoError(t *testing.T) {
	repoErr := errors.New("not found")
	svc := NewProfessionalService(&mockProfessionalRepo{
		getByUserID: func(_ context.Context, _ string) (*domain.Professional, error) { return nil, repoErr },
	}, nil)

	_, err := svc.GetByUserID(context.Background(), "user-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- Create ---

func TestProfessionalCreate_Success(t *testing.T) {
	input := ports.CreateProfessionalInput{
		UserID:            "user-1",
		LicenseNumber:     "LIC-001",
		Category:          domain.Nurse,
		YearsOfExperience: 5,
	}
	want := makeProfessional()

	svc := NewProfessionalService(&mockProfessionalRepo{
		create: func(_ context.Context, p *domain.Professional) (*domain.Professional, error) {
			if p.UserID != input.UserID {
				t.Errorf("UserID: got %s, want %s", p.UserID, input.UserID)
			}
			if p.LicenseNumber != input.LicenseNumber {
				t.Errorf("LicenseNumber: got %s, want %s", p.LicenseNumber, input.LicenseNumber)
			}
			if p.Category != input.Category {
				t.Errorf("Category: got %s, want %s", p.Category, input.Category)
			}
			if p.YearsOfExperience != input.YearsOfExperience {
				t.Errorf("YearsOfExperience: got %d, want %d", p.YearsOfExperience, input.YearsOfExperience)
			}
			return want, nil
		},
	}, nil)

	got, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestProfessionalCreate_WithResume(t *testing.T) {
	resume := "Experienced nurse"
	input := ports.CreateProfessionalInput{
		UserID:        "user-1",
		LicenseNumber: "LIC-001",
		Category:      domain.Nurse,
		Resume:        &resume,
	}
	svc := NewProfessionalService(&mockProfessionalRepo{
		create: func(_ context.Context, p *domain.Professional) (*domain.Professional, error) {
			if p.Resume != resume {
				t.Errorf("Resume: got %s, want %s", p.Resume, resume)
			}
			return p, nil
		},
	}, nil)

	if _, err := svc.Create(context.Background(), input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProfessionalCreate_EmptyUserID(t *testing.T) {
	svc := NewProfessionalService(&mockProfessionalRepo{}, nil)
	input := ports.CreateProfessionalInput{LicenseNumber: "LIC-001", Category: domain.Nurse}

	_, err := svc.Create(context.Background(), input)
	if !errors.Is(err, domain.ErrEmptyUserID) {
		t.Errorf("expected ErrEmptyUserID, got: %v", err)
	}
}

func TestProfessionalCreate_EmptyLicenseNumber(t *testing.T) {
	svc := NewProfessionalService(&mockProfessionalRepo{}, nil)
	input := ports.CreateProfessionalInput{UserID: "user-1", Category: domain.Nurse}

	_, err := svc.Create(context.Background(), input)
	if !errors.Is(err, domain.ErrEmptyLicenseNumber) {
		t.Errorf("expected ErrEmptyLicenseNumber, got: %v", err)
	}
}

func TestProfessionalCreate_InvalidCategory(t *testing.T) {
	svc := NewProfessionalService(&mockProfessionalRepo{}, nil)
	input := ports.CreateProfessionalInput{
		UserID:        "user-1",
		LicenseNumber: "LIC-001",
		Category:      "INVALID",
	}

	_, err := svc.Create(context.Background(), input)
	if !errors.Is(err, domain.ErrInvalidCategory) {
		t.Errorf("expected ErrInvalidCategory, got: %v", err)
	}
}

func TestProfessionalCreate_NegativeYears(t *testing.T) {
	svc := NewProfessionalService(&mockProfessionalRepo{}, nil)
	input := ports.CreateProfessionalInput{
		UserID:            "user-1",
		LicenseNumber:     "LIC-001",
		Category:          domain.Nurse,
		YearsOfExperience: -1,
	}

	_, err := svc.Create(context.Background(), input)
	if !errors.Is(err, domain.ErrNegativeYearsOfExp) {
		t.Errorf("expected ErrNegativeYearsOfExp, got: %v", err)
	}
}

func TestProfessionalCreate_RepoError(t *testing.T) {
	repoErr := errors.New("db error")
	svc := NewProfessionalService(&mockProfessionalRepo{
		create: func(_ context.Context, _ *domain.Professional) (*domain.Professional, error) { return nil, repoErr },
	}, nil)

	input := ports.CreateProfessionalInput{
		UserID:        "user-1",
		LicenseNumber: "LIC-001",
		Category:      domain.Nurse,
	}

	_, err := svc.Create(context.Background(), input)
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- Update ---

func TestProfessionalUpdate_AllFields(t *testing.T) {
	existing := makeProfessional()
	newLicense := "LIC-999"
	newCat := domain.Physiotherapist
	newYears := 10
	newVerified := true
	newResume := "Updated resume"

	input := ports.UpdateProfessionalInput{
		ID:                existing.ID,
		LicenseNumber:     &newLicense,
		Category:          &newCat,
		YearsOfExperience: &newYears,
		Verified:          &newVerified,
		Resume:            &newResume,
	}

	svc := NewProfessionalService(&mockProfessionalRepo{
		getByID: func(_ context.Context, _ string) (*domain.Professional, error) { return existing, nil },
		update: func(_ context.Context, p *domain.Professional) (*domain.Professional, error) {
			if p.LicenseNumber != newLicense {
				t.Errorf("LicenseNumber: got %s, want %s", p.LicenseNumber, newLicense)
			}
			if p.Category != newCat {
				t.Errorf("Category: got %s, want %s", p.Category, newCat)
			}
			if p.YearsOfExperience != newYears {
				t.Errorf("YearsOfExperience: got %d, want %d", p.YearsOfExperience, newYears)
			}
			if !p.Verified {
				t.Error("Verified should be true")
			}
			if p.Resume != newResume {
				t.Errorf("Resume: got %s, want %s", p.Resume, newResume)
			}
			return p, nil
		},
	}, nil)

	if _, err := svc.Update(context.Background(), input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProfessionalUpdate_PartialFields(t *testing.T) {
	existing := makeProfessional()
	newLicense := "LIC-002"

	svc := NewProfessionalService(&mockProfessionalRepo{
		getByID: func(_ context.Context, _ string) (*domain.Professional, error) { return existing, nil },
		update:  func(_ context.Context, p *domain.Professional) (*domain.Professional, error) { return p, nil },
	}, nil)

	got, err := svc.Update(context.Background(), ports.UpdateProfessionalInput{
		ID:            existing.ID,
		LicenseNumber: &newLicense,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.LicenseNumber != newLicense {
		t.Errorf("LicenseNumber: got %s, want %s", got.LicenseNumber, newLicense)
	}
	if got.Category != existing.Category {
		t.Errorf("Category should be unchanged: got %s, want %s", got.Category, existing.Category)
	}
}

func TestProfessionalUpdate_NotFound(t *testing.T) {
	repoErr := errors.New("not found")
	svc := NewProfessionalService(&mockProfessionalRepo{
		getByID: func(_ context.Context, _ string) (*domain.Professional, error) { return nil, repoErr },
	}, nil)

	_, err := svc.Update(context.Background(), ports.UpdateProfessionalInput{ID: "missing"})
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

func TestProfessionalUpdate_InvalidCategory(t *testing.T) {
	existing := makeProfessional()
	badCat := domain.Category("INVALID")

	svc := NewProfessionalService(&mockProfessionalRepo{
		getByID: func(_ context.Context, _ string) (*domain.Professional, error) { return existing, nil },
	}, nil)

	_, err := svc.Update(context.Background(), ports.UpdateProfessionalInput{
		ID:       existing.ID,
		Category: &badCat,
	})
	if !errors.Is(err, domain.ErrInvalidCategory) {
		t.Errorf("expected ErrInvalidCategory, got: %v", err)
	}
}

func TestProfessionalUpdate_EmptyLicenseNumber(t *testing.T) {
	existing := makeProfessional()
	empty := ""

	svc := NewProfessionalService(&mockProfessionalRepo{
		getByID: func(_ context.Context, _ string) (*domain.Professional, error) { return existing, nil },
	}, nil)

	_, err := svc.Update(context.Background(), ports.UpdateProfessionalInput{
		ID:            existing.ID,
		LicenseNumber: &empty,
	})
	if !errors.Is(err, domain.ErrEmptyLicenseNumber) {
		t.Errorf("expected ErrEmptyLicenseNumber, got: %v", err)
	}
}

func TestProfessionalUpdate_NegativeYears(t *testing.T) {
	existing := makeProfessional()
	neg := -1

	svc := NewProfessionalService(&mockProfessionalRepo{
		getByID: func(_ context.Context, _ string) (*domain.Professional, error) { return existing, nil },
	}, nil)

	_, err := svc.Update(context.Background(), ports.UpdateProfessionalInput{
		ID:                existing.ID,
		YearsOfExperience: &neg,
	})
	if !errors.Is(err, domain.ErrNegativeYearsOfExp) {
		t.Errorf("expected ErrNegativeYearsOfExp, got: %v", err)
	}
}

func TestProfessionalUpdate_RepoUpdateError(t *testing.T) {
	existing := makeProfessional()
	repoErr := errors.New("update failed")

	svc := NewProfessionalService(&mockProfessionalRepo{
		getByID: func(_ context.Context, _ string) (*domain.Professional, error) { return existing, nil },
		update:  func(_ context.Context, _ *domain.Professional) (*domain.Professional, error) { return nil, repoErr },
	}, nil)

	_, err := svc.Update(context.Background(), ports.UpdateProfessionalInput{
		ID:            existing.ID,
		LicenseNumber: ptr("LIC-002"),
	})
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- Delete ---

func TestProfessionalDelete_Success(t *testing.T) {
	svc := NewProfessionalService(&mockProfessionalRepo{
		delete: func(_ context.Context, id string) error {
			if id != "prof-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return nil
		},
	}, nil)

	if err := svc.Delete(context.Background(), "prof-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProfessionalDelete_RepoError(t *testing.T) {
	repoErr := errors.New("delete failed")
	svc := NewProfessionalService(&mockProfessionalRepo{
		delete: func(_ context.Context, _ string) error { return repoErr },
	}, nil)

	if err := svc.Delete(context.Background(), "prof-1"); !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- FindWithFilters ---

func TestProfessionalFindWithFilters_NoFilters(t *testing.T) {
	list := []*domain.Professional{makeProfessional(), makeProfessional()}
	svc := NewProfessionalService(&mockProfessionalRepo{
		findWithFilters: func(_ context.Context, f ports.ProfessionalFilters) ([]*domain.Professional, int64, error) {
			if f.City != nil || f.State != nil || len(f.DayOfWeek) != 0 || len(f.Shift) != 0 {
				t.Errorf("expected empty filters, got: %+v", f)
			}
			return list, int64(len(list)), nil
		},
	}, nil)

	got, err := svc.FindWithFilters(context.Background(), ports.ProfessionalFilters{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Items) != len(list) {
		t.Errorf("len: got %d, want %d", len(got.Items), len(list))
	}
}

func TestProfessionalFindWithFilters_WithFilters(t *testing.T) {
	city := "São Paulo"
	state := "SP"
	filters := ports.ProfessionalFilters{
		City:      &city,
		State:     &state,
		DayOfWeek: []string{"MONDAY", "WEDNESDAY"},
		Shift:     []string{"MORNING"},
	}
	list := []*domain.Professional{makeProfessional()}

	svc := NewProfessionalService(&mockProfessionalRepo{
		findWithFilters: func(_ context.Context, f ports.ProfessionalFilters) ([]*domain.Professional, int64, error) {
			if f.City == nil || *f.City != city {
				t.Errorf("City: got %v, want %s", f.City, city)
			}
			if f.State == nil || *f.State != state {
				t.Errorf("State: got %v, want %s", f.State, state)
			}
			return list, int64(len(list)), nil
		},
	}, nil)

	got, err := svc.FindWithFilters(context.Background(), filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Items) != 1 {
		t.Errorf("len: got %d, want 1", len(got.Items))
	}
}

func TestProfessionalFindWithFilters_EmptyResult(t *testing.T) {
	svc := NewProfessionalService(&mockProfessionalRepo{
		findWithFilters: func(_ context.Context, _ ports.ProfessionalFilters) ([]*domain.Professional, int64, error) {
			return []*domain.Professional{}, 0, nil
		},
	}, nil)

	got, err := svc.FindWithFilters(context.Background(), ports.ProfessionalFilters{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Items) != 0 {
		t.Errorf("expected empty slice, got %d items", len(got.Items))
	}
}

func TestProfessionalFindWithFilters_RepoError(t *testing.T) {
	repoErr := errors.New("db error")
	svc := NewProfessionalService(&mockProfessionalRepo{
		findWithFilters: func(_ context.Context, _ ports.ProfessionalFilters) ([]*domain.Professional, int64, error) {
			return nil, 0, repoErr
		},
	}, nil)

	_, err := svc.FindWithFilters(context.Background(), ports.ProfessionalFilters{})
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}
