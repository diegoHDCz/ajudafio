package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	userhttp "github.com/diegoHDCz/ajudafio/internal/user/adapters/http"
	"github.com/diegoHDCz/ajudafio/internal/user/domain"
	"github.com/diegoHDCz/ajudafio/internal/user/ports"
)

// --- Mock ---

type mockUserSvc struct {
	getByID        func(context.Context, string) (*domain.User, error)
	getByEmail     func(context.Context, string) (*domain.User, error)
	create         func(context.Context, ports.CreateUserInput) (*domain.User, error)
	update         func(context.Context, ports.UpdateUserInput) (*domain.User, error)
	deleteFn       func(context.Context, string) error
	updateRoleFn   func(context.Context, string, domain.Role) error
}

func (m *mockUserSvc) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return m.getByID(ctx, id)
}
func (m *mockUserSvc) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.getByEmail(ctx, email)
}
func (m *mockUserSvc) Create(ctx context.Context, input ports.CreateUserInput) (*domain.User, error) {
	return m.create(ctx, input)
}
func (m *mockUserSvc) Update(ctx context.Context, input ports.UpdateUserInput) (*domain.User, error) {
	return m.update(ctx, input)
}
func (m *mockUserSvc) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}
func (m *mockUserSvc) UpdateUserRole(ctx context.Context, id string, role domain.Role) error {
	if m.updateRoleFn != nil {
		return m.updateRoleFn(ctx, id, role)
	}
	return nil
}
func (m *mockUserSvc) UploadAvatar(_ context.Context, _ string, _ []byte, _ string) (*domain.User, error) {
	return nil, nil
}

func makeTestUser() *domain.User {
	return &domain.User{
		ID:    "user-1",
		Name:  "Alice",
		Email: "alice@example.com",
		Role:  domain.RoleClient,
	}
}

func newUserRouter(svc ports.UserService) http.Handler {
	return userhttp.NewRouter(userhttp.NewHandler(svc, shared.NewValidator(svc)))
}

// --- Me ---

func TestUserMe_WithClaims(t *testing.T) {
	claims := &authdomain.JWTClaims{Name: "Alice", Email: "alice@example.com"}
	router := newUserRouter(&mockUserSvc{})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	var resp struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Name != claims.Name || resp.Email != claims.Email {
		t.Errorf("response mismatch: got %+v, want {Name:%s Email:%s}", resp, claims.Name, claims.Email)
	}
}

func TestUserMe_NoClaims(t *testing.T) {
	router := newUserRouter(&mockUserSvc{})
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

// --- GetByID ---

func TestUserGetByID_Success(t *testing.T) {
	user := makeTestUser()
	svc := &mockUserSvc{
		getByID: func(_ context.Context, id string) (*domain.User, error) {
			if id != user.ID {
				t.Errorf("id: got %s, want %s", id, user.ID)
			}
			return user, nil
		},
	}
	router := newUserRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/"+string(user.ID), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	var resp struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID != string(user.ID) || resp.Email != user.Email {
		t.Errorf("response mismatch: got %+v", resp)
	}
}

func TestUserGetByID_NotFound(t *testing.T) {
	svc := &mockUserSvc{
		getByID: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
	}
	router := newUserRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}

// --- Create ---

func TestUserCreate_Success(t *testing.T) {
	user := makeTestUser()
	svc := &mockUserSvc{
		create: func(_ context.Context, input ports.CreateUserInput) (*domain.User, error) {
			if input.Email != user.Email || input.Name != user.Name {
				t.Errorf("input mismatch: %+v", input)
			}
			return user, nil
		},
	}
	body, _ := json.Marshal(map[string]string{"name": user.Name, "email": user.Email})
	router := newUserRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusCreated)
	}
	var resp struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID != string(user.ID) {
		t.Errorf("ID: got %s, want %s", resp.ID, user.ID)
	}
}

func TestUserCreate_InvalidJSON(t *testing.T) {
	router := newUserRouter(&mockUserSvc{})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestUserCreate_MissingName(t *testing.T) {
	router := newUserRouter(&mockUserSvc{})
	body, _ := json.Marshal(map[string]string{"email": "alice@example.com"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestUserCreate_MissingEmail(t *testing.T) {
	router := newUserRouter(&mockUserSvc{})
	body, _ := json.Marshal(map[string]string{"name": "Alice"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestUserCreate_ServiceError(t *testing.T) {
	svc := &mockUserSvc{
		create: func(_ context.Context, _ ports.CreateUserInput) (*domain.User, error) {
			return nil, errors.New("service error")
		},
	}
	body, _ := json.Marshal(map[string]string{"name": "Alice", "email": "alice@example.com"})
	router := newUserRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- Update ---

func TestUserUpdate_NoClaims(t *testing.T) {
	router := newUserRouter(&mockUserSvc{})
	body, _ := json.Marshal(map[string]string{"name": "X"})
	req := httptest.NewRequest(http.MethodPatch, "/user-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestUserUpdate_Success(t *testing.T) {
	user := makeTestUser()
	svc := &mockUserSvc{
		update: func(_ context.Context, input ports.UpdateUserInput) (*domain.User, error) {
			if input.ID != user.ID {
				t.Errorf("ID: got %s, want %s", input.ID, user.ID)
			}
			return user, nil
		},
	}
	body, _ := json.Marshal(map[string]string{"name": "Alice Updated"})
	router := newUserRouter(svc)
	req := httptest.NewRequest(http.MethodPatch, "/"+string(user.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.JWTClaims{
		Role: "admin",
	}))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestUserUpdate_InvalidJSON(t *testing.T) {
	router := newUserRouter(&mockUserSvc{})
	req := httptest.NewRequest(http.MethodPatch, "/user-1", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.JWTClaims{
		Role: "admin",
	}))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestUserUpdate_ServiceError(t *testing.T) {
	svc := &mockUserSvc{
		update: func(_ context.Context, _ ports.UpdateUserInput) (*domain.User, error) {
			return nil, errors.New("not found")
		},
	}
	body, _ := json.Marshal(map[string]string{"name": "X"})
	router := newUserRouter(svc)
	req := httptest.NewRequest(http.MethodPatch, "/user-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.JWTClaims{
		Role: "admin",
	}))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- Delete ---

func TestUserDelete_NoClaims(t *testing.T) {
	router := newUserRouter(&mockUserSvc{})
	req := httptest.NewRequest(http.MethodDelete, "/user-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestUserDelete_Success(t *testing.T) {
	svc := &mockUserSvc{
		deleteFn: func(_ context.Context, id string) error {
			if id != "user-1" {
				t.Errorf("id: got %s, want user-1", id)
			}
			return nil
		},
	}
	router := newUserRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/user-1", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.JWTClaims{
		Role: "admin",
	}))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestUserDelete_ServiceError(t *testing.T) {
	svc := &mockUserSvc{
		deleteFn: func(_ context.Context, _ string) error {
			return errors.New("delete failed")
		},
	}
	router := newUserRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/user-1", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.JWTClaims{
		Role: "admin",
	}))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
