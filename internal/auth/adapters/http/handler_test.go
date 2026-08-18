package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diegoHDCz/ajudafio/internal/auth"
	authhttp "github.com/diegoHDCz/ajudafio/internal/auth/adapters/http"
	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	"github.com/diegoHDCz/ajudafio/internal/auth/ports"
)

// --- Mock ---

type mockAuthSvc struct {
	register    func(context.Context, ports.RegisterInput) (*authdomain.TokenPair, error)
	login       func(context.Context, string, string) (*authdomain.TokenPair, error)
	refresh     func(context.Context, string) (*authdomain.TokenPair, error)
	logout      func(context.Context, string) error
	googleLogin func(context.Context, string) (*authdomain.TokenPair, error)
}

func (m *mockAuthSvc) Register(ctx context.Context, input ports.RegisterInput) (*authdomain.TokenPair, error) {
	return m.register(ctx, input)
}
func (m *mockAuthSvc) Login(ctx context.Context, email, password string) (*authdomain.TokenPair, error) {
	return m.login(ctx, email, password)
}
func (m *mockAuthSvc) Refresh(ctx context.Context, refreshToken string) (*authdomain.TokenPair, error) {
	return m.refresh(ctx, refreshToken)
}
func (m *mockAuthSvc) Logout(ctx context.Context, refreshToken string) error {
	return m.logout(ctx, refreshToken)
}
func (m *mockAuthSvc) GoogleLogin(ctx context.Context, idToken string) (*authdomain.TokenPair, error) {
	return m.googleLogin(ctx, idToken)
}

func newAuthRouter(svc ports.AuthService) http.Handler {
	return authhttp.NewRouter(authhttp.NewHandler(svc))
}

// --- GoogleLogin ---

func TestGoogleLogin_Success(t *testing.T) {
	svc := &mockAuthSvc{
		googleLogin: func(_ context.Context, idToken string) (*authdomain.TokenPair, error) {
			if idToken != "valid-token" {
				t.Errorf("id_token: got %s, want valid-token", idToken)
			}
			return &authdomain.TokenPair{AccessToken: "access", RefreshToken: "refresh", ExpiresIn: 3600}, nil
		},
	}
	body, _ := json.Marshal(map[string]string{"id_token": "valid-token"})
	req := httptest.NewRequest(http.MethodPost, "/google", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newAuthRouter(svc).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d, body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var resp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.AccessToken != "access" || resp.RefreshToken != "refresh" || resp.TokenType != "Bearer" {
		t.Errorf("response mismatch: got %+v", resp)
	}
}

func TestGoogleLogin_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/google", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newAuthRouter(&mockAuthSvc{}).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestGoogleLogin_MissingIDToken(t *testing.T) {
	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/google", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newAuthRouter(&mockAuthSvc{}).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestGoogleLogin_EmailNotVerified(t *testing.T) {
	svc := &mockAuthSvc{
		googleLogin: func(_ context.Context, _ string) (*authdomain.TokenPair, error) {
			return nil, auth.ErrGoogleEmailUnverified
		},
	}
	body, _ := json.Marshal(map[string]string{"id_token": "unverified-token"})
	req := httptest.NewRequest(http.MethodPost, "/google", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newAuthRouter(svc).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestGoogleLogin_InvalidToken(t *testing.T) {
	svc := &mockAuthSvc{
		googleLogin: func(_ context.Context, _ string) (*authdomain.TokenPair, error) {
			return nil, context.DeadlineExceeded
		},
	}
	body, _ := json.Marshal(map[string]string{"id_token": "bad-token"})
	req := httptest.NewRequest(http.MethodPost, "/google", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newAuthRouter(svc).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
