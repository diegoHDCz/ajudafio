package user

import (
	"context"
	"errors"
	"testing"

	"github.com/diegoHDCz/ajudafio/internal/user/domain"
	"github.com/diegoHDCz/ajudafio/internal/user/ports"
)

// mockUserRepo implements ports.UserRepository for testing.
type mockUserRepo struct {
	getByID    func(context.Context, string) (*domain.User, error)
	getByEmail func(context.Context, string) (*domain.User, error)
	create     func(context.Context, *domain.User) (*domain.User, error)
	update     func(context.Context, *domain.User) (*domain.User, error)
	delete     func(context.Context, string) error
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return m.getByID(ctx, id)
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.getByEmail(ctx, email)
}
func (m *mockUserRepo) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	return m.create(ctx, u)
}
func (m *mockUserRepo) Update(ctx context.Context, u *domain.User) (*domain.User, error) {
	return m.update(ctx, u)
}
func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	return m.delete(ctx, id)
}

func ptrString(s string) *string { return &s }
func ptrRole(r domain.Role) *domain.Role { return &r }

func makeUser() *domain.User {
	return &domain.User{
		ID:    "user-1",
		Name:  "Alice",
		Email: "alice@example.com",
		Role:  domain.RoleClient,
	}
}

// --- GetByID ---

func TestGetByID_Success(t *testing.T) {
	want := makeUser()
	svc := NewService(&mockUserRepo{
		getByID: func(_ context.Context, id string) (*domain.User, error) {
			if id != "user-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return want, nil
		},
	})

	got, err := svc.GetByID(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetByID_RepoError(t *testing.T) {
	repoErr := errors.New("db error")
	svc := NewService(&mockUserRepo{
		getByID: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, repoErr
		},
	})

	_, err := svc.GetByID(context.Background(), "user-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repoErr) {
		t.Errorf("expected wrapped repoErr, got: %v", err)
	}
}

// --- GetByEmail ---

func TestGetByEmail_Success(t *testing.T) {
	want := makeUser()
	svc := NewService(&mockUserRepo{
		getByEmail: func(_ context.Context, email string) (*domain.User, error) {
			if email != want.Email {
				t.Fatalf("unexpected email: %s", email)
			}
			return want, nil
		},
	})

	got, err := svc.GetByEmail(context.Background(), want.Email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetByEmail_RepoError(t *testing.T) {
	repoErr := errors.New("not found")
	svc := NewService(&mockUserRepo{
		getByEmail: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, repoErr
		},
	})

	_, err := svc.GetByEmail(context.Background(), "x@x.com")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- Create ---

func TestCreate_Success(t *testing.T) {
	input := ports.CreateUserInput{
		Email: "bob@example.com",
		Name:  "Bob",
		Role:  domain.RoleProfessional,
	}
	want := &domain.User{ID: "new-id", Email: input.Email, Name: input.Name, Role: input.Role}

	svc := NewService(&mockUserRepo{
		create: func(_ context.Context, u *domain.User) (*domain.User, error) {
			if u.Email != input.Email || u.Name != input.Name || u.Role != input.Role {
				t.Errorf("unexpected user passed to repo: %+v", u)
			}
			return want, nil
		},
	})

	got, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCreate_RepoError(t *testing.T) {
	repoErr := errors.New("constraint violation")
	svc := NewService(&mockUserRepo{
		create: func(_ context.Context, _ *domain.User) (*domain.User, error) {
			return nil, repoErr
		},
	})

	_, err := svc.Create(context.Background(), ports.CreateUserInput{Name: "X", Email: "x@x.com", Role: domain.RoleClient})
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- Update ---

func TestUpdate_AllFields(t *testing.T) {
	existing := makeUser()
	newName := "Alice Updated"
	newEmail := "new@example.com"
	newPhone := "123456"
	newRole := domain.RoleAdmin

	input := ports.UpdateUserInput{
		ID:    existing.ID,
		Name:  ptrString(newName),
		Email: ptrString(newEmail),
		Phone: ptrString(newPhone),
		Role:  ptrRole(newRole),
	}

	svc := NewService(&mockUserRepo{
		getByID: func(_ context.Context, id string) (*domain.User, error) {
			if id != existing.ID {
				t.Fatalf("unexpected id: %s", id)
			}
			return existing, nil
		},
		update: func(_ context.Context, u *domain.User) (*domain.User, error) {
			if u.Name != newName {
				t.Errorf("Name: got %s, want %s", u.Name, newName)
			}
			if u.Email != newEmail {
				t.Errorf("Email: got %s, want %s", u.Email, newEmail)
			}
			if u.Phone == nil || *u.Phone != newPhone {
				t.Errorf("Phone: got %v, want %s", u.Phone, newPhone)
			}
			if u.Role != newRole {
				t.Errorf("Role: got %s, want %s", u.Role, newRole)
			}
			return u, nil
		},
	})

	got, err := svc.Update(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != newName {
		t.Errorf("returned Name: got %s, want %s", got.Name, newName)
	}
}

func TestUpdate_PartialFields(t *testing.T) {
	existing := makeUser()
	newName := "Partial"

	input := ports.UpdateUserInput{ID: existing.ID, Name: ptrString(newName)}

	svc := NewService(&mockUserRepo{
		getByID: func(_ context.Context, _ string) (*domain.User, error) { return existing, nil },
		update:  func(_ context.Context, u *domain.User) (*domain.User, error) { return u, nil },
	})

	got, err := svc.Update(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != newName {
		t.Errorf("Name: got %s, want %s", got.Name, newName)
	}
	if got.Email != existing.Email {
		t.Errorf("Email should be unchanged: got %s, want %s", got.Email, existing.Email)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	repoErr := errors.New("not found")
	svc := NewService(&mockUserRepo{
		getByID: func(_ context.Context, _ string) (*domain.User, error) { return nil, repoErr },
	})

	_, err := svc.Update(context.Background(), ports.UpdateUserInput{ID: "missing"})
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

func TestUpdate_RepoUpdateError(t *testing.T) {
	existing := makeUser()
	repoErr := errors.New("update failed")

	svc := NewService(&mockUserRepo{
		getByID: func(_ context.Context, _ string) (*domain.User, error) { return existing, nil },
		update:  func(_ context.Context, _ *domain.User) (*domain.User, error) { return nil, repoErr },
	})

	_, err := svc.Update(context.Background(), ports.UpdateUserInput{ID: existing.ID, Name: ptrString("X")})
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}

// --- Delete ---

func TestDelete_Success(t *testing.T) {
	svc := NewService(&mockUserRepo{
		delete: func(_ context.Context, id string) error {
			if id != "user-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return nil
		},
	})

	if err := svc.Delete(context.Background(), "user-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDelete_RepoError(t *testing.T) {
	repoErr := errors.New("delete failed")
	svc := NewService(&mockUserRepo{
		delete: func(_ context.Context, _ string) error { return repoErr },
	})

	err := svc.Delete(context.Background(), "user-1")
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repoErr, got: %v", err)
	}
}
