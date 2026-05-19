package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	addrhttp "github.com/diegoHDCz/ajudafio/internal/address/adapters/http"
	"github.com/diegoHDCz/ajudafio/internal/address/domain"
	"github.com/diegoHDCz/ajudafio/internal/address/ports"
)

// --- Mock ---

type mockAddrSvc struct {
	getByID        func(context.Context, string) (*domain.Address, error)
	getByUserID    func(context.Context, string) ([]*domain.Address, error)
	getByContractID func(context.Context, string) ([]*domain.Address, error)
	create         func(context.Context, ports.CreateAddressInput) (*domain.Address, error)
	update         func(context.Context, ports.UpdateAddressInput) (*domain.Address, error)
	deleteFn       func(context.Context, string) error
}

func (m *mockAddrSvc) GetByID(ctx context.Context, id string) (*domain.Address, error) {
	return m.getByID(ctx, id)
}
func (m *mockAddrSvc) GetByUserID(ctx context.Context, userID string) ([]*domain.Address, error) {
	return m.getByUserID(ctx, userID)
}
func (m *mockAddrSvc) GetByContractID(ctx context.Context, contractID string) ([]*domain.Address, error) {
	return m.getByContractID(ctx, contractID)
}
func (m *mockAddrSvc) Create(ctx context.Context, input ports.CreateAddressInput) (*domain.Address, error) {
	return m.create(ctx, input)
}
func (m *mockAddrSvc) Update(ctx context.Context, input ports.UpdateAddressInput) (*domain.Address, error) {
	return m.update(ctx, input)
}
func (m *mockAddrSvc) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

func makeTestAddress() *domain.Address {
	contractID := "contract-1"
	return &domain.Address{
		ID:          "addr-1",
		UserID:      "user-1",
		ContractID:  &contractID,
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

func newAddrRouter(svc ports.AddressService) http.Handler {
	return addrhttp.NewRouter(addrhttp.NewAddressHandler(svc))
}

func adminClaims() *authdomain.JWTClaims {
	return &authdomain.JWTClaims{
		Sub:         "admin-user",
		RealmAccess: authdomain.RealmAccess{Roles: []string{"admin"}},
	}
}

// --- GetByID ---

func TestAddrGetByID_Success(t *testing.T) {
	a := makeTestAddress()
	svc := &mockAddrSvc{
		getByID: func(_ context.Context, id string) (*domain.Address, error) {
			if id != a.ID {
				t.Errorf("id: got %s, want %s", id, a.ID)
			}
			return a, nil
		},
	}
	router := newAddrRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/"+a.ID, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	var resp struct {
		ID   string `json:"id"`
		City string `json:"city"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID != a.ID {
		t.Errorf("ID: got %s, want %s", resp.ID, a.ID)
	}
	if resp.City != a.City {
		t.Errorf("City: got %s, want %s", resp.City, a.City)
	}
}

func TestAddrGetByID_NotFound(t *testing.T) {
	svc := &mockAddrSvc{
		getByID: func(_ context.Context, _ string) (*domain.Address, error) {
			return nil, errors.New("not found")
		},
	}
	router := newAddrRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}

// --- GetByUserID ---

func TestAddrGetByUserID_Success(t *testing.T) {
	list := []*domain.Address{makeTestAddress(), makeTestAddress()}
	svc := &mockAddrSvc{
		getByUserID: func(_ context.Context, userID string) ([]*domain.Address, error) {
			if userID != "user-1" {
				t.Errorf("userID: got %s, want user-1", userID)
			}
			return list, nil
		},
	}
	router := newAddrRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/user/user-1", nil)
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

func TestAddrGetByUserID_ServiceError(t *testing.T) {
	svc := &mockAddrSvc{
		getByUserID: func(_ context.Context, _ string) ([]*domain.Address, error) {
			return nil, errors.New("db error")
		},
	}
	router := newAddrRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/user/user-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- GetByContractID ---

func TestAddrGetByContractID_Success(t *testing.T) {
	list := []*domain.Address{makeTestAddress()}
	svc := &mockAddrSvc{
		getByContractID: func(_ context.Context, contractID string) ([]*domain.Address, error) {
			if contractID != "contract-1" {
				t.Errorf("contractID: got %s, want contract-1", contractID)
			}
			return list, nil
		},
	}
	router := newAddrRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/contract/contract-1", nil)
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

func TestAddrGetByContractID_ServiceError(t *testing.T) {
	svc := &mockAddrSvc{
		getByContractID: func(_ context.Context, _ string) ([]*domain.Address, error) {
			return nil, errors.New("db error")
		},
	}
	router := newAddrRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/contract/contract-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- Create ---

func validCreateBody() map[string]interface{} {
	return map[string]interface{}{
		"user_id":      "user-1",
		"zip_code":     "80000-000",
		"address_line": "Rua das Flores",
		"number":       "123",
		"district":     "Centro",
		"city":         "Curitiba",
		"state":        "PR",
	}
}

func TestAddrCreate_Success(t *testing.T) {
	a := makeTestAddress()
	svc := &mockAddrSvc{
		create: func(_ context.Context, input ports.CreateAddressInput) (*domain.Address, error) {
			if input.UserID != "user-1" {
				t.Errorf("UserID: got %s, want user-1", input.UserID)
			}
			if input.City != "Curitiba" {
				t.Errorf("City: got %s, want Curitiba", input.City)
			}
			return a, nil
		},
	}
	body, _ := json.Marshal(validCreateBody())
	router := newAddrRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusCreated)
	}
}

func TestAddrCreate_InvalidJSON(t *testing.T) {
	router := newAddrRouter(&mockAddrSvc{})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAddrCreate_MissingRequiredFields(t *testing.T) {
	cases := []struct {
		name string
		body map[string]interface{}
	}{
		{"missing user_id", map[string]interface{}{"zip_code": "80000-000", "address_line": "Rua A", "number": "1", "district": "D", "city": "C", "state": "PR"}},
		{"missing zip_code", map[string]interface{}{"user_id": "user-1", "address_line": "Rua A", "number": "1", "district": "D", "city": "C", "state": "PR"}},
		{"missing city", map[string]interface{}{"user_id": "user-1", "zip_code": "80000-000", "address_line": "Rua A", "number": "1", "district": "D", "state": "PR"}},
		{"missing state", map[string]interface{}{"user_id": "user-1", "zip_code": "80000-000", "address_line": "Rua A", "number": "1", "district": "D", "city": "C"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			router := newAddrRouter(&mockAddrSvc{})
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("%s: status: got %d, want %d", tc.name, rec.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestAddrCreate_ServiceError(t *testing.T) {
	svc := &mockAddrSvc{
		create: func(_ context.Context, _ ports.CreateAddressInput) (*domain.Address, error) {
			return nil, errors.New("db error")
		},
	}
	body, _ := json.Marshal(validCreateBody())
	router := newAddrRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

// --- Update ---

func TestAddrUpdate_Success(t *testing.T) {
	a := makeTestAddress()
	svc := &mockAddrSvc{
		update: func(_ context.Context, input ports.UpdateAddressInput) (*domain.Address, error) {
			if input.ID != a.ID {
				t.Errorf("ID: got %s, want %s", input.ID, a.ID)
			}
			return a, nil
		},
	}
	body, _ := json.Marshal(map[string]string{"city": "São Paulo"})
	router := newAddrRouter(svc)
	req := httptest.NewRequest(http.MethodPatch, "/"+a.ID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestAddrUpdate_InvalidJSON(t *testing.T) {
	router := newAddrRouter(&mockAddrSvc{})
	req := httptest.NewRequest(http.MethodPatch, "/addr-1", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAddrUpdate_ServiceError(t *testing.T) {
	svc := &mockAddrSvc{
		update: func(_ context.Context, _ ports.UpdateAddressInput) (*domain.Address, error) {
			return nil, errors.New("not found")
		},
	}
	body, _ := json.Marshal(map[string]string{"city": "São Paulo"})
	router := newAddrRouter(svc)
	req := httptest.NewRequest(http.MethodPatch, "/addr-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

// --- Delete ---

func TestAddrDelete_NoClaims(t *testing.T) {
	router := newAddrRouter(&mockAddrSvc{})
	req := httptest.NewRequest(http.MethodDelete, "/addr-1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAddrDelete_NotAdmin(t *testing.T) {
	claims := &authdomain.JWTClaims{Sub: "user-1"}
	router := newAddrRouter(&mockAddrSvc{})
	req := httptest.NewRequest(http.MethodDelete, "/addr-1", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestAddrDelete_AdminSuccess(t *testing.T) {
	svc := &mockAddrSvc{
		deleteFn: func(_ context.Context, id string) error {
			if id != "addr-1" {
				t.Errorf("id: got %s, want addr-1", id)
			}
			return nil
		},
	}
	router := newAddrRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/addr-1", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), adminClaims()))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestAddrDelete_ServiceError(t *testing.T) {
	svc := &mockAddrSvc{
		deleteFn: func(_ context.Context, _ string) error { return errors.New("delete failed") },
	}
	router := newAddrRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/addr-1", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), adminClaims()))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
