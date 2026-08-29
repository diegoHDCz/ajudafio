package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	authhttp "github.com/diegoHDCz/ajudafio/internal/auth/adapters/http"
	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	profiledomain "github.com/diegoHDCz/ajudafio/internal/profile/domain"
	profileports "github.com/diegoHDCz/ajudafio/internal/profile/ports"
	userdomain "github.com/diegoHDCz/ajudafio/internal/user/domain"
	userports "github.com/diegoHDCz/ajudafio/internal/user/ports"
)

// --- Mocks ---

type mockUserSvc struct {
	getByID           func(context.Context, string) (*userdomain.User, error)
	ensureProvisioned func(context.Context, userports.ProvisionUserInput) (*userdomain.User, error)
	updateOnboarding  func(context.Context, string, userdomain.OnboardingStatus) (*userdomain.User, error)
	listRoles         func(context.Context, string) ([]userdomain.Role, error)
	updateUserRole    func(context.Context, string, userdomain.Role) error
}

func (m *mockUserSvc) GetByID(ctx context.Context, id string) (*userdomain.User, error) {
	return m.getByID(ctx, id)
}
func (m *mockUserSvc) GetByEmail(_ context.Context, _ string) (*userdomain.User, error) {
	return nil, nil
}
func (m *mockUserSvc) GetByAuthUserID(_ context.Context, _ string) (*userdomain.User, error) {
	return nil, nil
}
func (m *mockUserSvc) Create(_ context.Context, _ userports.CreateUserInput) (*userdomain.User, error) {
	return nil, nil
}
func (m *mockUserSvc) Update(_ context.Context, _ userports.UpdateUserInput) (*userdomain.User, error) {
	return nil, nil
}
func (m *mockUserSvc) Delete(_ context.Context, _ string) error { return nil }
func (m *mockUserSvc) UpdateUserRole(ctx context.Context, id string, role userdomain.Role) error {
	if m.updateUserRole != nil {
		return m.updateUserRole(ctx, id, role)
	}
	return nil
}
func (m *mockUserSvc) UploadAvatar(_ context.Context, _ string, _ []byte, _ string) (*userdomain.User, error) {
	return nil, nil
}
func (m *mockUserSvc) EnsureProvisioned(ctx context.Context, input userports.ProvisionUserInput) (*userdomain.User, error) {
	return m.ensureProvisioned(ctx, input)
}
func (m *mockUserSvc) UpdateOnboardingStatus(ctx context.Context, id string, status userdomain.OnboardingStatus) (*userdomain.User, error) {
	return m.updateOnboarding(ctx, id, status)
}
func (m *mockUserSvc) ListRoles(ctx context.Context, id string) ([]userdomain.Role, error) {
	return m.listRoles(ctx, id)
}

type mockProfileSvc struct{}

func (m *mockProfileSvc) CreateFamilyProfile(_ context.Context, userID string) (*profiledomain.FamilyProfile, error) {
	return &profiledomain.FamilyProfile{ID: "fam-1", UserID: userID}, nil
}
func (m *mockProfileSvc) GetFamilyProfileByUserID(_ context.Context, _ string) (*profiledomain.FamilyProfile, error) {
	return nil, nil
}
func (m *mockProfileSvc) CreateFinancialProfile(_ context.Context, userID string) (*profiledomain.FinancialProfile, error) {
	return &profiledomain.FinancialProfile{ID: "fin-1", UserID: userID}, nil
}
func (m *mockProfileSvc) GetFinancialProfileByUserID(_ context.Context, _ string) (*profiledomain.FinancialProfile, error) {
	return nil, nil
}

var _ profileports.ProfileService = (*mockProfileSvc)(nil)

type mockAuditSvc struct{}

func (m *mockAuditSvc) Log(_ context.Context, _ *string, _ string, _ map[string]any) {}

func newAuthRouter(userSvc *mockUserSvc) http.Handler {
	return authhttp.NewRouter(authhttp.NewHandler(userSvc, &mockProfileSvc{}, &mockAuditSvc{}))
}

// --- Me ---

func TestMe_ExistingUser(t *testing.T) {
	userSvc := &mockUserSvc{
		getByID: func(_ context.Context, id string) (*userdomain.User, error) {
			if id != "user-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return &userdomain.User{ID: "user-1", Name: "Alice", Email: "alice@example.com", Role: userdomain.RoleFamilyClient}, nil
		},
	}
	h := authhttp.NewHandler(userSvc, &mockProfileSvc{}, &mockAuditSvc{})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.AuthenticatedUser{
		AuthUserID: "auth-1", UserID: "user-1", Email: "alice@example.com",
	}))
	rec := httptest.NewRecorder()
	h.Me(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d, body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var resp struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID != "user-1" || resp.Name != "Alice" {
		t.Errorf("response mismatch: got %+v", resp)
	}
}

func TestMe_FirstAccessProvisions(t *testing.T) {
	provisioned := false
	userSvc := &mockUserSvc{
		ensureProvisioned: func(_ context.Context, input userports.ProvisionUserInput) (*userdomain.User, error) {
			provisioned = true
			if input.AuthUserID != "auth-2" || input.Email != "new@example.com" {
				t.Errorf("unexpected provision input: %+v", input)
			}
			return &userdomain.User{ID: "user-2", Name: "new@example.com", Email: "new@example.com", Role: userdomain.RoleFamilyClient}, nil
		},
	}
	h := authhttp.NewHandler(userSvc, &mockProfileSvc{}, &mockAuditSvc{})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.AuthenticatedUser{
		AuthUserID: "auth-2", Email: "new@example.com",
	}))
	rec := httptest.NewRecorder()
	h.Me(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d, body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if !provisioned {
		t.Error("expected EnsureProvisioned to be called for first access")
	}
}

func TestMe_NoClaims(t *testing.T) {
	h := authhttp.NewHandler(&mockUserSvc{}, &mockProfileSvc{}, &mockAuditSvc{})
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()
	h.Me(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

// --- MeRoles / MePermissions ---

func TestMeRoles_Success(t *testing.T) {
	userSvc := &mockUserSvc{
		listRoles: func(_ context.Context, _ string) ([]userdomain.Role, error) {
			return []userdomain.Role{userdomain.RoleFamilyClient}, nil
		},
	}
	h := authhttp.NewHandler(userSvc, &mockProfileSvc{}, &mockAuditSvc{})

	req := httptest.NewRequest(http.MethodGet, "/me/roles", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.AuthenticatedUser{UserID: "user-1"}))
	rec := httptest.NewRecorder()
	h.MeRoles(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestMeRoles_Unprovisioned(t *testing.T) {
	h := authhttp.NewHandler(&mockUserSvc{}, &mockProfileSvc{}, &mockAuditSvc{})
	req := httptest.NewRequest(http.MethodGet, "/me/roles", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.AuthenticatedUser{}))
	rec := httptest.NewRecorder()
	h.MeRoles(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestMePermissions_Success(t *testing.T) {
	userSvc := &mockUserSvc{
		listRoles: func(_ context.Context, _ string) ([]userdomain.Role, error) {
			return []userdomain.Role{userdomain.RolePlatformAdmin}, nil
		},
	}
	h := authhttp.NewHandler(userSvc, &mockProfileSvc{}, &mockAuditSvc{})

	req := httptest.NewRequest(http.MethodGet, "/me/permissions", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.AuthenticatedUser{UserID: "user-1"}))
	rec := httptest.NewRecorder()
	h.MePermissions(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	var resp struct {
		Permissions []string `json:"permissions"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Permissions) == 0 {
		t.Error("expected PLATFORM_ADMIN to have permissions")
	}
}

// --- CompleteProfile ---

func TestCompleteProfile_FamilyClient(t *testing.T) {
	roleUpdated := false
	userSvc := &mockUserSvc{
		updateUserRole: func(_ context.Context, id string, role userdomain.Role) error {
			roleUpdated = true
			if id != "user-1" || role != userdomain.RoleFamilyClient {
				t.Errorf("unexpected update: id=%s role=%s", id, role)
			}
			return nil
		},
	}
	router := newAuthRouter(userSvc)

	body, _ := json.Marshal(map[string]string{"role": "FAMILY_CLIENT"})
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewReader(body))
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.AuthenticatedUser{UserID: "user-1"}))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status: got %d, want %d, body: %s", rec.Code, http.StatusNoContent, rec.Body.String())
	}
	if !roleUpdated {
		t.Error("expected role to be updated")
	}
}

func TestCompleteProfile_HealthCareproviderRejected(t *testing.T) {
	router := newAuthRouter(&mockUserSvc{})

	body, _ := json.Marshal(map[string]string{"role": "HEALTH_CAREPROVIDER"})
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewReader(body))
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.AuthenticatedUser{UserID: "user-1"}))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestCompleteProfile_InvalidRole(t *testing.T) {
	router := newAuthRouter(&mockUserSvc{})

	body, _ := json.Marshal(map[string]string{"role": "NOT_A_ROLE"})
	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewReader(body))
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.AuthenticatedUser{UserID: "user-1"}))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// --- Onboarding ---

func TestUpdateOnboarding_Success(t *testing.T) {
	userSvc := &mockUserSvc{
		updateOnboarding: func(_ context.Context, id string, status userdomain.OnboardingStatus) (*userdomain.User, error) {
			if id != "user-1" || status != userdomain.OnboardingInProgress {
				t.Errorf("unexpected update: id=%s status=%s", id, status)
			}
			return &userdomain.User{ID: id, OnboardingStatus: status}, nil
		},
	}
	h := authhttp.NewHandler(userSvc, &mockProfileSvc{}, &mockAuditSvc{})

	body, _ := json.Marshal(map[string]string{"status": "IN_PROGRESS"})
	req := httptest.NewRequest(http.MethodPost, "/onboarding", bytes.NewReader(body))
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.AuthenticatedUser{UserID: "user-1"}))
	rec := httptest.NewRecorder()
	h.UpdateOnboarding(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d, body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestOnboardingStatus_Success(t *testing.T) {
	userSvc := &mockUserSvc{
		getByID: func(_ context.Context, _ string) (*userdomain.User, error) {
			return &userdomain.User{ID: "user-1", OnboardingStatus: userdomain.OnboardingCompleted}, nil
		},
	}
	h := authhttp.NewHandler(userSvc, &mockProfileSvc{}, &mockAuditSvc{})

	req := httptest.NewRequest(http.MethodGet, "/onboarding/status", nil)
	req = req.WithContext(authmiddleware.WithClaims(req.Context(), &authdomain.AuthenticatedUser{UserID: "user-1"}))
	rec := httptest.NewRecorder()
	h.OnboardingStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	var resp onboardingStatusResp
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != "COMPLETED" {
		t.Errorf("status: got %s, want COMPLETED", resp.Status)
	}
}

type onboardingStatusResp struct {
	Status string `json:"status"`
}
