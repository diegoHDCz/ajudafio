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
	profhttp "github.com/diegoHDCz/ajudafio/internal/professional/adapters/http"
	"github.com/diegoHDCz/ajudafio/internal/professional/domain"
	"github.com/diegoHDCz/ajudafio/internal/professional/ports"
)

// --- Mock ---

type mockProfSvc struct {
	getByID         func(context.Context, string) (*domain.Professional, error)
	getByUserID     func(context.Context, string) (*domain.Professional, error)
	create          func(context.Context, ports.CreateProfessionalInput) (*domain.Professional, error)
	update          func(context.Context, ports.UpdateProfessionalInput) (*domain.Professional, error)
	deleteFn        func(context.Context, string) error
	findWithFilters func(context.Context, ports.ProfessionalFilters) ([]*domain.Professional, error)
}

func (m *mockProfSvc) GetByID(ctx context.Context, id string) (*domain.Professional, error) {
	return m.getByID(ctx, id)
}
func (m *mockProfSvc) GetByUserID(ctx context.Context, userID string) (*domain.Professional, error) {
	return m.getByUserID(ctx, userID)
}
func (m *mockProfSvc) Create(ctx context.Context, input ports.CreateProfessionalInput) (*domain.Professional, error) {
	return m.create(ctx, input)
}
func (m *mockProfSvc) Update(ctx context.Context, input ports.UpdateProfessionalInput) (*domain.Professional, error) {
	return m.update(ctx, input)
}
func (m *mockProfSvc) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}
func (m *mockProfSvc) FindWithFilters(ctx context.Context, filters ports.ProfessionalFilters) ([]*domain.Professional, error) {
	return m.findWithFilters(ctx, filters)
}

func makeTestProfessional() *domain.Professional {
	return &domain.Professional{
		ID:                "prof-1",
		UserID:            "user-1",
		LicenseNumber:     "LIC-001",
		Category:          domain.Nurse,
		YearsOfExperience: 5,
	}
}

func newProfRouter(svc ports.ProfessionalService) http.Handler {
	return profhttp.NewRouter(profhttp.NewProfessionalHandler(svc))
}

// --- FindWithFilters ---

func TestProfFindWithFilters_NoFilters(t *testing.T) {
	list := []*domain.Professional{makeTestProfessional()}
	svc := &mockProfSvc{
		findWithFilters: func(_ context.Context, f ports.ProfessionalFilters) ([]*domain.Professional, error) {
			if f.City != nil || f.State != nil || len(f.DayOfWeek) != 0 || len(f.Shift) != 0 {
				t.Errorf("expected empty filters, got %+v", f)
			}
			return list, nil
		},
	}
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	var resp []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp) != 1 {
		t.Errorf("len: got %d, want 1", len(resp))
	}
}

func TestProfFindWithFilters_WithFilters(t *testing.T) {
	svc := &mockProfSvc{
		findWithFilters: func(_ context.Context, f ports.ProfessionalFilters) ([]*domain.Professional, error) {
			if f.City == nil || *f.City != "curitiba" {
				t.Errorf("City: got %v, want curitiba", f.City)
			}
			if f.State == nil || *f.State != "PR" {
				t.Errorf("State: got %v, want PR", f.State)
			}
			if len(f.DayOfWeek) != 2 {
				t.Errorf("DayOfWeek: got %d items, want 2", len(f.DayOfWeek))
			}
			if len(f.Shift) != 1 || f.Shift[0] != "MORNING" {
				t.Errorf("Shift: got %v, want [MORNING]", f.Shift)
			}
			return []*domain.Professional{makeTestProfessional()}, nil
		},
	}
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodGet,
		"/?city=curitiba&state=PR&day_of_week=MONDAY&day_of_week=WEDNESDAY&shift=MORNING", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	var resp []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp) != 1 {
		t.Errorf("len: got %d, want 1", len(resp))
	}
}

func TestProfFindWithFilters_ServiceError(t *testing.T) {
	svc := &mockProfSvc{
		findWithFilters: func(_ context.Context, _ ports.ProfessionalFilters) ([]*domain.Professional, error) {
			return nil, errors.New("db error")
		},
	}
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- GetByUserID ---

func TestProfGetByUserID_Success(t *testing.T) {
	p := makeTestProfessional()
	svc := &mockProfSvc{
		getByUserID: func(_ context.Context, id string) (*domain.Professional, error) {
			if id != p.UserID {
				t.Errorf("userID: got %s, want %s", id, p.UserID)
			}
			return p, nil
		},
	}
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/user/"+p.UserID, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	var resp struct {
		ID     string `json:"id"`
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.UserID != p.UserID {
		t.Errorf("UserID: got %s, want %s", resp.UserID, p.UserID)
	}
}

func TestProfGetByUserID_NotFound(t *testing.T) {
	svc := &mockProfSvc{
		getByUserID: func(_ context.Context, _ string) (*domain.Professional, error) {
			return nil, errors.New("not found")
		},
	}
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/user/unknown", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}

// --- GetByID ---

func TestProfGetByID_Success(t *testing.T) {
	p := makeTestProfessional()
	svc := &mockProfSvc{
		getByID: func(_ context.Context, id string) (*domain.Professional, error) {
			if id != p.ID {
				t.Errorf("id: got %s, want %s", id, p.ID)
			}
			return p, nil
		},
	}
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/"+p.ID, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	var resp struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID != p.ID {
		t.Errorf("ID: got %s, want %s", resp.ID, p.ID)
	}
}

func TestProfGetByID_NotFound(t *testing.T) {
	svc := &mockProfSvc{
		getByID: func(_ context.Context, _ string) (*domain.Professional, error) {
			return nil, errors.New("not found")
		},
	}
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}

// --- Create ---

func TestProfCreate_Success(t *testing.T) {
	p := makeTestProfessional()
	svc := &mockProfSvc{
		create: func(_ context.Context, input ports.CreateProfessionalInput) (*domain.Professional, error) {
			if input.UserID != p.UserID || input.LicenseNumber != p.LicenseNumber {
				t.Errorf("input mismatch: %+v", input)
			}
			return p, nil
		},
	}
	body, _ := json.Marshal(map[string]interface{}{
		"user_id":        p.UserID,
		"license_number": p.LicenseNumber,
		"category":       string(p.Category),
	})
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusCreated)
	}
}

func TestProfCreate_InvalidJSON(t *testing.T) {
	router := newProfRouter(&mockProfSvc{})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestProfCreate_MissingUserID(t *testing.T) {
	router := newProfRouter(&mockProfSvc{})
	body, _ := json.Marshal(map[string]string{"license_number": "LIC-001"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestProfCreate_MissingLicenseNumber(t *testing.T) {
	router := newProfRouter(&mockProfSvc{})
	body, _ := json.Marshal(map[string]string{"user_id": "user-1"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestProfCreate_ServiceValidationError(t *testing.T) {
	svc := &mockProfSvc{
		create: func(_ context.Context, _ ports.CreateProfessionalInput) (*domain.Professional, error) {
			return nil, domain.ErrInvalidCategory
		},
	}
	body, _ := json.Marshal(map[string]string{
		"user_id": "user-1", "license_number": "LIC-001", "category": "INVALID",
	})
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

// --- Update ---

func TestProfUpdate_Success(t *testing.T) {
	p := makeTestProfessional()
	svc := &mockProfSvc{
		update: func(_ context.Context, input ports.UpdateProfessionalInput) (*domain.Professional, error) {
			if input.ID != p.ID {
				t.Errorf("ID: got %s, want %s", input.ID, p.ID)
			}
			return p, nil
		},
	}
	body, _ := json.Marshal(map[string]interface{}{"license_number": "LIC-002"})
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodPatch, "/"+p.ID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestProfUpdate_InvalidJSON(t *testing.T) {
	router := newProfRouter(&mockProfSvc{})
	req := httptest.NewRequest(http.MethodPatch, "/prof-1", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestProfUpdate_ServiceError(t *testing.T) {
	svc := &mockProfSvc{
		update: func(_ context.Context, _ ports.UpdateProfessionalInput) (*domain.Professional, error) {
			return nil, errors.New("not found")
		},
	}
	body, _ := json.Marshal(map[string]string{"license_number": "LIC-002"})
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodPatch, "/prof-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

// --- Delete ---

func TestProfDelete_NoClaims(t *testing.T) {
	router := newProfRouter(&mockProfSvc{})
	req := httptest.NewRequest(http.MethodDelete, "/prof-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestProfDelete_ProfessionalNotFound(t *testing.T) {
	claims := &authdomain.JWTClaims{Sub: "user-1"}
	svc := &mockProfSvc{
		getByID: func(_ context.Context, _ string) (*domain.Professional, error) {
			return nil, errors.New("not found")
		},
	}
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/prof-1", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestProfDelete_Forbidden(t *testing.T) {
	p := makeTestProfessional() // UserID = "user-1"
	claims := &authdomain.JWTClaims{Sub: "other-user"}
	svc := &mockProfSvc{
		getByID: func(_ context.Context, _ string) (*domain.Professional, error) { return p, nil },
	}
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/"+p.ID, nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestProfDelete_OwnerCanDelete(t *testing.T) {
	p := makeTestProfessional() // UserID = "user-1"
	claims := &authdomain.JWTClaims{Sub: p.UserID}
	svc := &mockProfSvc{
		getByID:  func(_ context.Context, _ string) (*domain.Professional, error) { return p, nil },
		deleteFn: func(_ context.Context, id string) error {
			if id != p.ID {
				t.Errorf("id: got %s, want %s", id, p.ID)
			}
			return nil
		},
	}
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/"+p.ID, nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestProfDelete_AdminCanDelete(t *testing.T) {
	p := makeTestProfessional() // UserID = "user-1"
	claims := &authdomain.JWTClaims{
		Sub:         "other-user",
		RealmAccess: authdomain.RealmAccess{Roles: []string{"admin"}},
	}
	svc := &mockProfSvc{
		getByID:  func(_ context.Context, _ string) (*domain.Professional, error) { return p, nil },
		deleteFn: func(_ context.Context, _ string) error { return nil },
	}
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/"+p.ID, nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestProfDelete_ServiceError(t *testing.T) {
	p := makeTestProfessional()
	claims := &authdomain.JWTClaims{Sub: p.UserID}
	svc := &mockProfSvc{
		getByID:  func(_ context.Context, _ string) (*domain.Professional, error) { return p, nil },
		deleteFn: func(_ context.Context, _ string) error { return errors.New("delete failed") },
	}
	router := newProfRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/"+p.ID, nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
