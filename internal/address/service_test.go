package address

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/diegoHDCz/ajudafio/internal/address/domain"
	"github.com/diegoHDCz/ajudafio/internal/address/ports"
)

// --- Mock ---

type mockAddressRepo struct {
	createAddress        func(*domain.Address) error
	getAddressByID       func(string) (*domain.Address, error)
	updateAddress        func(*domain.Address) error
	deleteAddress        func(string) error
	getAddressesByUserID func(string) ([]*domain.Address, error)
	getAllAddresses       func() ([]*domain.Address, error)
	getAddressesByCity   func(string) ([]*domain.Address, error)
}

func (m *mockAddressRepo) CreateAddress(a *domain.Address) error {
	return m.createAddress(a)
}
func (m *mockAddressRepo) GetAddressByID(id string) (*domain.Address, error) {
	return m.getAddressByID(id)
}
func (m *mockAddressRepo) UpdateAddress(a *domain.Address) error {
	return m.updateAddress(a)
}
func (m *mockAddressRepo) DeleteAddress(id string) error {
	return m.deleteAddress(id)
}
func (m *mockAddressRepo) GetAddressesByUserID(userID string) ([]*domain.Address, error) {
	return m.getAddressesByUserID(userID)
}
func (m *mockAddressRepo) GetAllAddresses() ([]*domain.Address, error) {
	return m.getAllAddresses()
}
func (m *mockAddressRepo) GetAddressesByCity(city string) ([]*domain.Address, error) {
	return m.getAddressesByCity(city)
}

func ptr[T any](v T) *T { return &v }

func makeAddress() *domain.Address {
	return &domain.Address{
		ID:          "addr-1",
		UserID:      "user-1",
		ZipCode:     "80000-000",
		AddressLine: "Rua das Flores",
		Number:      "123",
		District:    "Centro",
		City:        "Curitiba",
		State:       "PR",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// --- GetByID ---

func TestAddressGetByID_Success(t *testing.T) {
	want := makeAddress()
	svc := NewAddressService(&mockAddressRepo{
		getAddressByID: func(id string) (*domain.Address, error) {
			if id != want.ID {
				t.Fatalf("unexpected id: %s", id)
			}
			return want, nil
		},
	})

	got, err := svc.GetByID(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAddressGetByID_RepoError(t *testing.T) {
	repoErr := errors.New("not found")
	svc := NewAddressService(&mockAddressRepo{
		getAddressByID: func(_ string) (*domain.Address, error) { return nil, repoErr },
	})

	_, err := svc.GetByID(context.Background(), "addr-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- GetByUserID ---

func TestAddressGetByUserID_Success(t *testing.T) {
	list := []*domain.Address{makeAddress()}
	svc := NewAddressService(&mockAddressRepo{
		getAddressesByUserID: func(userID string) ([]*domain.Address, error) {
			if userID != "user-1" {
				t.Fatalf("unexpected userID: %s", userID)
			}
			return list, nil
		},
	})

	got, err := svc.GetByUserID(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len: got %d, want 1", len(got))
	}
}

func TestAddressGetByUserID_RepoError(t *testing.T) {
	repoErr := errors.New("db error")
	svc := NewAddressService(&mockAddressRepo{
		getAddressesByUserID: func(_ string) ([]*domain.Address, error) { return nil, repoErr },
	})

	_, err := svc.GetByUserID(context.Background(), "user-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- Create ---

func TestAddressCreate_Success(t *testing.T) {
	input := ports.CreateAddressInput{
		UserID:      "user-1",
		ZipCode:     "80000-000",
		AddressLine: "Rua das Flores",
		Number:      "123",
		District:    "Centro",
		City:        "Curitiba",
		State:       "PR",
	}

	svc := NewAddressService(&mockAddressRepo{
		createAddress: func(a *domain.Address) error {
			if a.UserID != input.UserID {
				t.Errorf("UserID: got %s, want %s", a.UserID, input.UserID)
			}
			if a.ZipCode != input.ZipCode {
				t.Errorf("ZipCode: got %s, want %s", a.ZipCode, input.ZipCode)
			}
			if a.City != input.City {
				t.Errorf("City: got %s, want %s", a.City, input.City)
			}
			if a.ID == "" {
				t.Error("expected non-empty ID")
			}
			return nil
		},
	})

	got, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.UserID != input.UserID {
		t.Errorf("UserID: got %s, want %s", got.UserID, input.UserID)
	}
	if got.ID == "" {
		t.Error("expected returned address to have an ID")
	}
}

func TestAddressCreate_WithOptionalFields(t *testing.T) {
	complement := "Apto 42"
	reference := "Próximo ao mercado"
	input := ports.CreateAddressInput{
		UserID:      "user-1",
		ZipCode:     "80000-000",
		AddressLine: "Rua das Flores",
		Number:      "123",
		District:    "Centro",
		City:        "Curitiba",
		State:       "PR",
		Complement:  &complement,
		Reference:   &reference,
	}

	svc := NewAddressService(&mockAddressRepo{
		createAddress: func(a *domain.Address) error {
			if a.Complement == nil || *a.Complement != complement {
				t.Errorf("Complement: got %v, want %s", a.Complement, complement)
			}
			if a.Reference == nil || *a.Reference != reference {
				t.Errorf("Reference: got %v, want %s", a.Reference, reference)
			}
			return nil
		},
	})

	if _, err := svc.Create(context.Background(), input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddressCreate_RepoError(t *testing.T) {
	repoErr := errors.New("db error")
	svc := NewAddressService(&mockAddressRepo{
		createAddress: func(_ *domain.Address) error { return repoErr },
	})

	_, err := svc.Create(context.Background(), ports.CreateAddressInput{UserID: "user-1"})
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- Update ---

func TestAddressUpdate_AllFields(t *testing.T) {
	existing := makeAddress()

	input := ports.UpdateAddressInput{
		ID:          existing.ID,
		ZipCode:     ptr("90000-000"),
		AddressLine: ptr("Av. Paulista"),
		Number:      ptr("999"),
		Complement:  ptr("Sala 1"),
		District:    ptr("Bela Vista"),
		City:        ptr("São Paulo"),
		State:       ptr("SP"),
		Reference:   ptr("Em frente ao MASP"),
	}

	svc := NewAddressService(&mockAddressRepo{
		getAddressByID: func(_ string) (*domain.Address, error) { return existing, nil },
		updateAddress: func(a *domain.Address) error {
			if a.ZipCode != *input.ZipCode {
				t.Errorf("ZipCode: got %s, want %s", a.ZipCode, *input.ZipCode)
			}
			if a.City != *input.City {
				t.Errorf("City: got %s, want %s", a.City, *input.City)
			}
			if a.State != *input.State {
				t.Errorf("State: got %s, want %s", a.State, *input.State)
			}
			return nil
		},
	})

	got, err := svc.Update(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.City != *input.City {
		t.Errorf("City: got %s, want %s", got.City, *input.City)
	}
}

func TestAddressUpdate_PartialFields(t *testing.T) {
	existing := makeAddress()
	originalCity := existing.City

	svc := NewAddressService(&mockAddressRepo{
		getAddressByID: func(_ string) (*domain.Address, error) { return existing, nil },
		updateAddress:  func(_ *domain.Address) error { return nil },
	})

	got, err := svc.Update(context.Background(), ports.UpdateAddressInput{
		ID:      existing.ID,
		ZipCode: ptr("99999-999"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ZipCode != "99999-999" {
		t.Errorf("ZipCode: got %s, want 99999-999", got.ZipCode)
	}
	if got.City != originalCity {
		t.Errorf("City should be unchanged: got %s, want %s", got.City, originalCity)
	}
}

func TestAddressUpdate_NotFound(t *testing.T) {
	repoErr := errors.New("address not found")
	svc := NewAddressService(&mockAddressRepo{
		getAddressByID: func(_ string) (*domain.Address, error) { return nil, repoErr },
	})

	_, err := svc.Update(context.Background(), ports.UpdateAddressInput{ID: "missing"})
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

func TestAddressUpdate_RepoUpdateError(t *testing.T) {
	existing := makeAddress()
	repoErr := errors.New("update failed")

	svc := NewAddressService(&mockAddressRepo{
		getAddressByID: func(_ string) (*domain.Address, error) { return existing, nil },
		updateAddress:  func(_ *domain.Address) error { return repoErr },
	})

	_, err := svc.Update(context.Background(), ports.UpdateAddressInput{
		ID:   existing.ID,
		City: ptr("São Paulo"),
	})
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- Delete ---

func TestAddressDelete_Success(t *testing.T) {
	svc := NewAddressService(&mockAddressRepo{
		deleteAddress: func(id string) error {
			if id != "addr-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return nil
		},
	})

	if err := svc.Delete(context.Background(), "addr-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddressDelete_RepoError(t *testing.T) {
	repoErr := errors.New("delete failed")
	svc := NewAddressService(&mockAddressRepo{
		deleteAddress: func(_ string) error { return repoErr },
	})

	if err := svc.Delete(context.Background(), "addr-1"); !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}
