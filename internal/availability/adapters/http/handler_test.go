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
	getByID                func(context.Context, string) (*domain.Availability, error)
	getByProfessionalID    func(context.Context, string) ([]*domain.Availability, error)
	create                 func(context.Context, *domain.Availability) (*domain.Availability, error)
	update                 func(context.Context, *domain.Availability) (*domain.Availability, error)
	delete                 func(context.Context, string) error
	deleteByProfessionalID func(context.Context, string) error
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
func (r *stubAvailRepo) DeleteByProfessionalID(ctx context.Context, id string) error {
	if r.deleteByProfessionalID != nil {
		return r.deleteByProfessionalID(ctx, id)
	}
	return nil
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
func (s *stubProfSvc) FindWithFilters(_ context.Context, _ profports.ProfessionalFilters) (*profports.ProfessionalPage, error) {
	return nil, errors.New("not implemented")
}

// --- Stub user service (for validator) ---

type stubUserSvcAvail struct {
	getByEmail func(context.Context, string) (*userdomain.User, error)
}

func (s *stubUserSvcAvail) GetByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	if s.getByEmail != nil {
		return s.getByEmail(ctx, email)
	}
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
func (s *stubUserSvcAvail) UploadAvatar(_ context.Context, _ string, _ []byte, _ string) (*userdomain.User, error) {
	return nil, errors.New("not implemented")
}

func ptrShift(s shared.Shift) *shared.Shift { return &s }

func makeTestAvailability() *domain.Availability {
	return &domain.Availability{
		ID:             "avail-1",
		ProfessionalID: "prof-1",
		DayOfWeek:      shared.Monday,
		Shift:          ptrShift(shared.ShiftMorning),
	}
}

func adminClaims() *authdomain.JWTClaims {
	return &authdomain.JWTClaims{Role: "admin"}
}

func newAvailRouter(repo *stubAvailRepo) http.Handler {
	return newAvailRouterFull(repo, &stubProfSvc{}, &stubUserSvcAvail{})
}

func newAvailRouterFull(repo *stubAvailRepo, profSvc *stubProfSvc, userSvc *stubUserSvcAvail) http.Handler {
	svc := avail.NewAvailabilityService(repo)
	validator := shared.NewValidator(userSvc)
	h := availhttp.NewAvailabilityHandler(svc, validator, profSvc)
	return availhttp.NewAvailabilityRouter(h)
}

func ownerClaims() *authdomain.JWTClaims {
	return &authdomain.JWTClaims{Email: "owner@test.com"}
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
		deleteByProfessionalID: func(_ context.Context, _ string) error { return nil },
		create: func(_ context.Context, a *domain.Availability) (*domain.Availability, error) {
			if a.ProfessionalID != "prof-1" {
				t.Errorf("ProfessionalID: got %s, want prof-1", a.ProfessionalID)
			}
			if a.DayOfWeek != shared.Monday {
				t.Errorf("DayOfWeek: got %s, want MONDAY", a.DayOfWeek)
			}
			return created, nil
		},
	}
	body, _ := json.Marshal(map[string]interface{}{
		"professional_id": "prof-1",
		"availabilities": []map[string]interface{}{
			{"day_of_week": "MONDAY", "shift": "MORNING"},
		},
	})
	router := newAvailRouter(repo)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d — body: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	var resp []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp) != 1 {
		t.Errorf("len: got %d, want 1", len(resp))
	}
}

func TestAvailHandlerCreate_ShiftResolvesToHours(t *testing.T) {
	repo := &stubAvailRepo{
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
	}
	body, _ := json.Marshal(map[string]interface{}{
		"professional_id": "prof-1",
		"availabilities": []map[string]interface{}{
			{"day_of_week": "MONDAY", "shift": "MORNING"},
		},
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

func TestAvailHandlerCreate_MultipleIntervalsPerDay(t *testing.T) {
	createCount := 0
	start1, end1 := "10:00", "12:00"
	start2, end2 := "16:00", "20:00"
	repo := &stubAvailRepo{
		deleteByProfessionalID: func(_ context.Context, _ string) error { return nil },
		create: func(_ context.Context, a *domain.Availability) (*domain.Availability, error) {
			createCount++
			return a, nil
		},
	}
	body, _ := json.Marshal(map[string]interface{}{
		"professional_id": "prof-1",
		"availabilities": []map[string]interface{}{
			{"day_of_week": "WEDNESDAY", "start_hour": start1, "end_hour": end1},
			{"day_of_week": "WEDNESDAY", "start_hour": start2, "end_hour": end2},
		},
	})
	router := newAvailRouter(repo)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d — %s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if createCount != 2 {
		t.Errorf("create called %d times, want 2", createCount)
	}
}

func TestAvailHandlerCreate_OverlapRejected(t *testing.T) {
	repo := &stubAvailRepo{}
	body, _ := json.Marshal(map[string]interface{}{
		"professional_id": "prof-1",
		"availabilities": []map[string]interface{}{
			{"day_of_week": "WEDNESDAY", "start_hour": "10:00", "end_hour": "14:00"},
			{"day_of_week": "WEDNESDAY", "start_hour": "12:00", "end_hour": "16:00"},
		},
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
	body, _ := json.Marshal(map[string]interface{}{
		"availabilities": []map[string]interface{}{
			{"day_of_week": "MONDAY", "shift": "MORNING"},
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAvailHandlerCreate_EmptyAvailabilities(t *testing.T) {
	router := newAvailRouter(&stubAvailRepo{})
	body, _ := json.Marshal(map[string]interface{}{
		"professional_id": "prof-1",
		"availabilities":  []interface{}{},
	})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAvailHandlerCreate_CustomMissingHours(t *testing.T) {
	router := newAvailRouter(&stubAvailRepo{})
	body, _ := json.Marshal(map[string]interface{}{
		"professional_id": "prof-1",
		"availabilities": []map[string]interface{}{
			{"day_of_week": "MONDAY"},
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

func TestAvailHandlerCreate_ServiceError(t *testing.T) {
	repo := &stubAvailRepo{
		deleteByProfessionalID: func(_ context.Context, _ string) error {
			return errors.New("delete failed")
		},
	}
	body, _ := json.Marshal(map[string]interface{}{
		"professional_id": "prof-1",
		"availabilities": []map[string]interface{}{
			{"day_of_week": "MONDAY", "shift": "MORNING"},
		},
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
	body, _ := json.Marshal(map[string]interface{}{"day_of_week": "TUESDAY"})
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
	body, _ := json.Marshal(map[string]interface{}{"day_of_week": "TUESDAY", "shift": "AFTERNOON"})
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
	body, _ := json.Marshal(map[string]interface{}{"day_of_week": "TUESDAY"})
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

func ownerRouter(repo *stubAvailRepo) http.Handler {
	profSvc := &stubProfSvc{
		getByID: func(_ context.Context, _ string) (*profdomain.Professional, error) {
			return &profdomain.Professional{UserID: "user-1"}, nil
		},
	}
	userSvc := &stubUserSvcAvail{
		getByEmail: func(_ context.Context, _ string) (*userdomain.User, error) {
			return &userdomain.User{ID: "user-1"}, nil
		},
	}
	repo.getByID = func(_ context.Context, _ string) (*domain.Availability, error) {
		return makeTestAvailability(), nil
	}
	return newAvailRouterFull(repo, profSvc, userSvc)
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
	router := ownerRouter(repo)
	req := httptest.NewRequest(http.MethodDelete, "/avail-1", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), ownerClaims()))
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
	router := ownerRouter(repo)
	req := httptest.NewRequest(http.MethodDelete, "/avail-1", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), ownerClaims()))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
