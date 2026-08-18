package auth

import (
	"context"
	"errors"
	"testing"

	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	"github.com/diegoHDCz/ajudafio/internal/auth/ports"
	userdomain "github.com/diegoHDCz/ajudafio/internal/user/domain"
	userports "github.com/diegoHDCz/ajudafio/internal/user/ports"
	"github.com/google/uuid"
)

// --- Mocks ---

type mockAuthRepo struct {
	getUserAuthByEmail    func(context.Context, string) (*authdomain.UserAuth, error)
	setPasswordHash       func(context.Context, string, string) error
	createRefreshToken    func(context.Context, *authdomain.RefreshToken) (*authdomain.RefreshToken, error)
	getRefreshTokenByHash func(context.Context, string) (*authdomain.RefreshToken, error)
	revokeRefreshToken    func(context.Context, string) error
	revokeAllForUser      func(context.Context, string) error
	deleteExpiredTokens   func(context.Context) error
	getIdentityByProvider func(context.Context, string, string) (*authdomain.Identity, error)
	createIdentity        func(context.Context, uuid.UUID, string, string, string) (*authdomain.Identity, error)
}

func (m *mockAuthRepo) GetUserAuthByEmail(ctx context.Context, email string) (*authdomain.UserAuth, error) {
	return m.getUserAuthByEmail(ctx, email)
}
func (m *mockAuthRepo) SetPasswordHash(ctx context.Context, userID, hash string) error {
	return m.setPasswordHash(ctx, userID, hash)
}
func (m *mockAuthRepo) CreateRefreshToken(ctx context.Context, token *authdomain.RefreshToken) (*authdomain.RefreshToken, error) {
	if m.createRefreshToken != nil {
		return m.createRefreshToken(ctx, token)
	}
	return token, nil
}
func (m *mockAuthRepo) GetRefreshTokenByHash(ctx context.Context, hash string) (*authdomain.RefreshToken, error) {
	return m.getRefreshTokenByHash(ctx, hash)
}
func (m *mockAuthRepo) RevokeRefreshToken(ctx context.Context, hash string) error {
	return m.revokeRefreshToken(ctx, hash)
}
func (m *mockAuthRepo) RevokeAllRefreshTokensForUser(ctx context.Context, userID string) error {
	return m.revokeAllForUser(ctx, userID)
}
func (m *mockAuthRepo) DeleteExpiredRefreshTokens(ctx context.Context) error {
	return m.deleteExpiredTokens(ctx)
}
func (m *mockAuthRepo) GetIdentityByProvider(ctx context.Context, provider, providerUserID string) (*authdomain.Identity, error) {
	if m.getIdentityByProvider != nil {
		return m.getIdentityByProvider(ctx, provider, providerUserID)
	}
	return nil, nil
}
func (m *mockAuthRepo) CreateIdentity(ctx context.Context, userID uuid.UUID, provider, providerUserID, email string) (*authdomain.Identity, error) {
	if m.createIdentity != nil {
		return m.createIdentity(ctx, userID, provider, providerUserID, email)
	}
	return &authdomain.Identity{UserID: userID.String(), Provider: provider, ProviderUserID: providerUserID, Email: email}, nil
}

type mockUserSvc struct {
	getByID    func(context.Context, string) (*userdomain.User, error)
	getByEmail func(context.Context, string) (*userdomain.User, error)
	create     func(context.Context, userports.CreateUserInput) (*userdomain.User, error)
}

func (m *mockUserSvc) GetByID(ctx context.Context, id string) (*userdomain.User, error) {
	return m.getByID(ctx, id)
}
func (m *mockUserSvc) GetByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	return m.getByEmail(ctx, email)
}
func (m *mockUserSvc) Create(ctx context.Context, input userports.CreateUserInput) (*userdomain.User, error) {
	return m.create(ctx, input)
}
func (m *mockUserSvc) Update(ctx context.Context, input userports.UpdateUserInput) (*userdomain.User, error) {
	return nil, nil
}
func (m *mockUserSvc) Delete(ctx context.Context, id string) error { return nil }
func (m *mockUserSvc) UpdateUserRole(ctx context.Context, id string, role userdomain.Role) error {
	return nil
}
func (m *mockUserSvc) UploadAvatar(ctx context.Context, userID string, fileData []byte, contentType string) (*userdomain.User, error) {
	return nil, nil
}

type mockGoogleVerifier struct {
	verify func(context.Context, string) (*authdomain.GoogleClaims, error)
}

func (m *mockGoogleVerifier) Verify(ctx context.Context, idToken string) (*authdomain.GoogleClaims, error) {
	return m.verify(ctx, idToken)
}

func newTestAuthService(repo ports.AuthRepository, userSvc userports.UserService, verifier GoogleTokenVerifier) ports.AuthService {
	return NewService(repo, userSvc, []byte("test-secret"), verifier)
}

// --- GoogleLogin ---

func TestGoogleLogin_NewUser(t *testing.T) {
	googleUserID := uuid.New()
	var createdIdentityUserID uuid.UUID

	repo := &mockAuthRepo{
		getIdentityByProvider: func(_ context.Context, provider, providerUserID string) (*authdomain.Identity, error) {
			if provider != "google" || providerUserID != "google-sub-1" {
				t.Errorf("unexpected lookup: %s/%s", provider, providerUserID)
			}
			return nil, nil
		},
		createIdentity: func(_ context.Context, userID uuid.UUID, provider, providerUserID, email string) (*authdomain.Identity, error) {
			createdIdentityUserID = userID
			return &authdomain.Identity{UserID: userID.String(), Provider: provider, ProviderUserID: providerUserID, Email: email}, nil
		},
	}
	userSvc := &mockUserSvc{
		getByEmail: func(_ context.Context, _ string) (*userdomain.User, error) {
			return nil, errors.New("user not found")
		},
		create: func(_ context.Context, input userports.CreateUserInput) (*userdomain.User, error) {
			if input.Email != "new@example.com" || input.Role != userdomain.RoleClient {
				t.Errorf("unexpected create input: %+v", input)
			}
			return &userdomain.User{ID: googleUserID.String(), Email: input.Email, Name: input.Name, Role: input.Role}, nil
		},
	}
	verifier := &mockGoogleVerifier{
		verify: func(_ context.Context, idToken string) (*authdomain.GoogleClaims, error) {
			return &authdomain.GoogleClaims{Subject: "google-sub-1", Email: "new@example.com", EmailVerified: true, Name: "New User"}, nil
		},
	}

	svc := newTestAuthService(repo, userSvc, verifier)
	pair, err := svc.GoogleLogin(context.Background(), "some-id-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Errorf("expected non-empty token pair, got %+v", pair)
	}
	if createdIdentityUserID != googleUserID {
		t.Errorf("identity linked to wrong user: got %s, want %s", createdIdentityUserID, googleUserID)
	}
}

func TestGoogleLogin_ExistingIdentity(t *testing.T) {
	repo := &mockAuthRepo{
		getIdentityByProvider: func(_ context.Context, _ string, providerUserID string) (*authdomain.Identity, error) {
			return &authdomain.Identity{UserID: "user-1", Provider: "google", ProviderUserID: providerUserID}, nil
		},
	}
	userSvc := &mockUserSvc{
		getByID: func(_ context.Context, id string) (*userdomain.User, error) {
			if id != "user-1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return &userdomain.User{ID: "user-1", Email: "existing@example.com", Name: "Existing", Role: userdomain.RoleClient}, nil
		},
	}
	verifier := &mockGoogleVerifier{
		verify: func(_ context.Context, _ string) (*authdomain.GoogleClaims, error) {
			return &authdomain.GoogleClaims{Subject: "google-sub-2", Email: "existing@example.com", EmailVerified: true}, nil
		},
	}

	svc := newTestAuthService(repo, userSvc, verifier)
	pair, err := svc.GoogleLogin(context.Background(), "some-id-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pair.AccessToken == "" {
		t.Errorf("expected access token")
	}
}

func TestGoogleLogin_LinksExistingPasswordUser(t *testing.T) {
	existingID := uuid.New()
	var identityLinkedTo uuid.UUID

	repo := &mockAuthRepo{
		getIdentityByProvider: func(_ context.Context, _ string, _ string) (*authdomain.Identity, error) {
			return nil, nil
		},
		createIdentity: func(_ context.Context, userID uuid.UUID, provider, providerUserID, email string) (*authdomain.Identity, error) {
			identityLinkedTo = userID
			return &authdomain.Identity{UserID: userID.String(), Provider: provider, ProviderUserID: providerUserID, Email: email}, nil
		},
	}
	userSvc := &mockUserSvc{
		getByEmail: func(_ context.Context, email string) (*userdomain.User, error) {
			return &userdomain.User{ID: existingID.String(), Email: email, Name: "Password User", Role: userdomain.RoleClient}, nil
		},
		create: func(_ context.Context, _ userports.CreateUserInput) (*userdomain.User, error) {
			t.Fatal("Create should not be called when a user with the same email already exists")
			return nil, nil
		},
	}
	verifier := &mockGoogleVerifier{
		verify: func(_ context.Context, _ string) (*authdomain.GoogleClaims, error) {
			return &authdomain.GoogleClaims{Subject: "google-sub-3", Email: "password@example.com", EmailVerified: true}, nil
		},
	}

	svc := newTestAuthService(repo, userSvc, verifier)
	pair, err := svc.GoogleLogin(context.Background(), "some-id-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pair.AccessToken == "" {
		t.Errorf("expected access token")
	}
	if identityLinkedTo != existingID {
		t.Errorf("identity linked to wrong user: got %s, want %s", identityLinkedTo, existingID)
	}
}

func TestGoogleLogin_EmailNotVerified(t *testing.T) {
	verifier := &mockGoogleVerifier{
		verify: func(_ context.Context, _ string) (*authdomain.GoogleClaims, error) {
			return &authdomain.GoogleClaims{Subject: "google-sub-4", Email: "unverified@example.com", EmailVerified: false}, nil
		},
	}
	svc := newTestAuthService(&mockAuthRepo{}, &mockUserSvc{}, verifier)

	_, err := svc.GoogleLogin(context.Background(), "some-id-token")
	if !errors.Is(err, ErrGoogleEmailUnverified) {
		t.Errorf("expected ErrGoogleEmailUnverified, got: %v", err)
	}
}

func TestGoogleLogin_InvalidToken(t *testing.T) {
	verifyErr := errors.New("token expired")
	verifier := &mockGoogleVerifier{
		verify: func(_ context.Context, _ string) (*authdomain.GoogleClaims, error) {
			return nil, verifyErr
		},
	}
	svc := newTestAuthService(&mockAuthRepo{}, &mockUserSvc{}, verifier)

	_, err := svc.GoogleLogin(context.Background(), "bad-token")
	if !errors.Is(err, verifyErr) {
		t.Errorf("expected wrapped verifyErr, got: %v", err)
	}
}
