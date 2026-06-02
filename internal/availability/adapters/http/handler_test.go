package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	avail "github.com/diegoHDCz/ajudafio/internal/availability"
	availhttp "github.com/diegoHDCz/ajudafio/internal/availability/adapters/http"
	"github.com/diegoHDCz/ajudafio/internal/availability/domain"
	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	profdomain "github.com/diegoHDCz/ajudafio/internal/professional/domain"
	profports "github.com/diegoHDCz/ajudafio/internal/professional/ports"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	userdomain "github.com/diegoHDCz/ajudafio/internal/user/domain"
	userports "github.com/diegoHDCz/ajudafio/internal/user/ports"
)

// --- Stub repository ---

type stubAvailRepo struct {
	getByID             func(context.Context, string) (*domain.Availability, error)
	getByProfessionalID func(context.Context, string) ([]*domain.Availability, error)
	create              func(context.Context, *domain.Availability) (*domain.Availability, error)
	update              func(context.Context, *domain.Availability) (*domain.Availability, error)
	delete              func(context.Context, string) error
}

func (r *stubAvailRepo) GetByID(ctx context.Context, id string) (*domain.Availability, error) {
	if r.getByID != nil {
		return r.getByID(ctx, id)
	}
	return nil, errors.New("not found")
}
func (r *stubAvailRepo) GetByProfessionalID(ctx context.Context, id string) ([]*domain.Availability, error) {
	return r.getByProfessionalID(ctx, id)
}
func (r *stubAvailRepo) Create(ctx context.Context, a *domain.Availability) (*domain.Availability, error) {
	return r.create(ctx, a)
}
func (r *stubAvailRepo) Update(ctx context.Context, a *domain.Availability) (*domain.Availability, error) {
	return r.update(ctx, a)
}
func (r *stubAvailRepo) Delete(ctx context.Context, id string) error {
	return r.delete(ctx, id)
}

// --- Stub professional service ---

type stubProfSvc struct {
	getByID func(context.Context, string) (*profdomain.Professional, error)
}

func (s *stubProfSvc) GetByID(ctx context.Context, id string) (*profdomain.Professional, error) {
	if s.getByID != nil {
		return s.getByID(ctx, id)
	}
	return nil, errors.New("not found")
}
func (s *stubProfSvc) GetByUserID(_ context.Context, _ string) (*profdomain.Professional, error) {
	return nil, errors.New("not implemented")
}
func (s *stubProfSvc) Create(_ context.Context, _ profports.CreateProfessionalInput) (*profdomain.Professional, error) {
	return nil, errors.New("not implemented")
}
func (s *stubProfSvc) Update(_ context.Context, _ profports.UpdateProfessionalInput) (*profdomain.Professional, error) {
	return nil, errors.New("not implemented")
}
func (s *stubProfSvc) Delete(_ context.Context, _ string) error { return errors.New("not implemented") }
func (s *stubProfSvc) FindWithFilters(_ context.Context, _ profports.ProfessionalFilters) ([]*profdomain.Professional, error) {
	return nil, errors.New("not implemented")
}

// --- Stub user service (for validator) ---

type stubUserSvcAvail struct{}

func (s *stubUserSvcAvail) GetByEmail(_ context.Context, _ string) (*userdomain.User, error) {
	return nil, errors.New("not found")
}
func (s *stubUserSvcAvail) GetByID(_ context.Context, _ string) (*userdomain.User, error) {
	return nil, errors.New("not implemented")
}
func (s *stubUserSvcAvail) Create(_ context.Context, _ userports.CreateUserInput) (*userdomain.User, error) {
	return nil, errors.New("not implemented")
}
func (s *stubUserSvcAvail) Update(_ context.Context, _ userports.UpdateUserInput) (*userdomain.User, error) {
	return nil, errors.New("not implemented")
}
func (s *stubUserSvcAvail) Delete(_ context.Context, _ string) error {
	return errors.New("not implemented")
}
func (s *stubUserSvcAvail) UpdateUserRole(_ context.Context, _ string, _ userdomain.Role) error {
	return errors.New("not implemented")
}

func makeTestAvailability() *domain.Availability {
	shifts := []shared.Shift{shared.ShiftMorning}
	return &domain.Availability{
		ID:             "avail-1",
		ProfessionalID: "prof-1",
		DayOfWeek:      []shared.WeekDay{shared.Monday},
		Shift:          &shifts,
	}
}

func adminClaims() *authdomain.JWTClaims {
	return &authdomain.JWTClaims{Role: "admin"}
}

func newAvailRouter(repo *stubAvailRepo) http.Handler {
	svc := avail.NewAvailabilityService(repo)
	validator := shared.NewValidator(&stubUserSvcAvail{})
	h := availhttp.NewAvailabilityHandler(svc, validator, &stubProfSvc{})
	return availhttp.NewAvailabilityRouter(h)
}

// --- GetByProfessionalID ---

func TestAvailHandlerGetByProfessionalID_Success(t *testing.T) {
	list := []*domain.Availability{makeTestAvailability(), makeTestAvailability()}
	repo := &stubAvailRepo{
		getByProfessionalID: func(_ context.Context, id string) ([]*domain.Availability, error) {
			if id != "prof-1" {
				t.Errorf("id: got %s, want prof-1", id)
			}
			return list, nil
		},
	}
	router := newAvailRouter(repo)
	req := httptest.NewRequest(http.MethodGet, "/professional/prof-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	var resp []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp) != 2 {
		t.Errorf("len: got %d, want 2", len(resp))
	}
}

func TestAvailHandlerGetByProfessionalID_Empty(t *testing.T) {
	repo := &stubAvailRepo{
		getByProfessionalID: func(_ context.Context, _ string) ([]*domain.Availability, error) {
			return []*domain.Availability{}, nil
		},
	}
	router := newAvailRouter(repo)
	req := httptest.NewRequest(http.MethodGet, "/professional/prof-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	var resp []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp) != 0 {
		t.Errorf("expected empty slice, got %d items", len(resp))
	}
}

func TestAvailHandlerGetByProfessionalID_ServiceError(t *testing.T) {
	repo := &stubAvailRepo{
		getByProfessionalID: func(_ context.Context, _ string) ([]*domain.Availability, error) {
			return nil, errors.New("db error")
		},
	}
	router := newAvailRouter(repo)
	req := httptest.NewRequest(http.MethodGet, "/professional/prof-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- Create ---

func TestAvailHandlerCreate_Success(t *testing.T) {
	created := makeTestAvailability()
	created.ID = "new-avail-id"
	repo := &stubAvailRepo{
		create: func(_ context.Context, a *domain.Availability) (*domain.Availability, error) {
			if a.ProfessionalID != "prof-1" {
				t.Errorf("ProfessionalID: got %s, want prof-1", a.ProfessionalID)
			}
			if len(a.DayOfWeek) != 1 {
				t.Errorf("DayOfWeek len: got %d, want 1", len(a.DayOfWeek))
			}
			return created, nil
		},
	}
	body, _ := json.Marshal(map[string]interface{}{
		"professional_id": "prof-1",
		"day_of_week":     []string{"MONDAY"},
	})
	router := newAvailRouter(repo)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusCreated)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["id"] != created.ID {
		t.Errorf("id: got %v, want %s", resp["id"], created.ID)
	}
}

func TestAvailHandlerCreate_WithShift(t *testing.T) {
	created := makeTestAvailability()
	repo := &stubAvailRepo{
		create: func(_ context.Context, a *domain.Availability) (*domain.Availability, error) {
			if a.Shift == nil || len(*a.Shift) != 1 {
				t.Error("expected shift to be set")
			}
			return created, nil
		},
	}
	body, _ := json.Marshal(map[string]interface{}{
		"professional_id": "prof-1",
		"day_of_week":     []string{"MONDAY"},
		"shift":           []string{"MORNING"},
	})
	router := newAvailRouter(repo)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusCreated)
	}
}

func TestAvailHandlerCreate_InvalidJSON(t *testing.T) {
	router := newAvailRouter(&stubAvailRepo{})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAvailHandlerCreate_MissingProfessionalID(t *testing.T) {
	router := newAvailRouter(&stubAvailRepo{})
	body, _ := json.Marshal(map[string]interface{}{"day_of_week": []string{"MONDAY"}})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAvailHandlerCreate_MissingDayOfWeek(t *testing.T) {
	router := newAvailRouter(&stubAvailRepo{})
	body, _ := json.Marshal(map[string]interface{}{
		"professional_id": "prof-1",
		"day_of_week":     []string{},
	})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAvailHandlerCreate_ServiceError(t *testing.T) {
	repo := &stubAvailRepo{
		create: func(_ context.Context, _ *domain.Availability) (*domain.Availability, error) {
			return nil, errors.New("create failed")
		},
	}
	body, _ := json.Marshal(map[string]interface{}{
		"professional_id": "prof-1",
		"day_of_week":     []string{"MONDAY"},
	})
	router := newAvailRouter(repo)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

// --- Update ---

func TestAvailHandlerUpdate_NoClaims(t *testing.T) {
	router := newAvailRouter(&stubAvailRepo{})
	body, _ := json.Marshal(map[string]interface{}{"day_of_week": []string{"TUESDAY"}})
	req := httptest.NewRequest(http.MethodPatch, "/avail-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAvailHandlerUpdate_Success(t *testing.T) {
	a := makeTestAvailability()
	repo := &stubAvailRepo{
		update: func(_ context.Context, av *domain.Availability) (*domain.Availability, error) {
			if av.ID != "avail-1" {
				t.Errorf("ID: got %s, want avail-1", av.ID)
			}
			return a, nil
		},
	}
	body, _ := json.Marshal(map[string]interface{}{"day_of_week": []string{"TUESDAY", "THURSDAY"}})
	router := newAvailRouter(repo)
	req := httptest.NewRequest(http.MethodPatch, "/avail-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), adminClaims()))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestAvailHandlerUpdate_InvalidJSON(t *testing.T) {
	router := newAvailRouter(&stubAvailRepo{})
	req := httptest.NewRequest(http.MethodPatch, "/avail-1", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), adminClaims()))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAvailHandlerUpdate_ServiceError(t *testing.T) {
	repo := &stubAvailRepo{
		update: func(_ context.Context, _ *domain.Availability) (*domain.Availability, error) {
			return nil, errors.New("update failed")
		},
	}
	body, _ := json.Marshal(map[string]interface{}{"day_of_week": []string{"TUESDAY"}})
	router := newAvailRouter(repo)
	req := httptest.NewRequest(http.MethodPatch, "/avail-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), adminClaims()))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

// --- Delete ---

func TestAvailHandlerDelete_NoClaims(t *testing.T) {
	router := newAvailRouter(&stubAvailRepo{})
	req := httptest.NewRequest(http.MethodDelete, "/avail-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAvailHandlerDelete_Success(t *testing.T) {
	repo := &stubAvailRepo{
		delete: func(_ context.Context, id string) error {
			if id != "avail-1" {
				t.Errorf("id: got %s, want avail-1", id)
			}
			return nil
		},
	}
	router := newAvailRouter(repo)
	req := httptest.NewRequest(http.MethodDelete, "/avail-1", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), adminClaims()))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestAvailHandlerDelete_ServiceError(t *testing.T) {
	repo := &stubAvailRepo{
		delete: func(_ context.Context, _ string) error { return errors.New("delete failed") },
	}
	router := newAvailRouter(repo)
	req := httptest.NewRequest(http.MethodDelete, "/avail-1", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), adminClaims()))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
