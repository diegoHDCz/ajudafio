package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/diegoHDCz/ajudafio/internal/auth/domain"
	userports "github.com/diegoHDCz/ajudafio/internal/user/ports"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const claimsKey contextKey = "claims"

// supabaseClaims is the raw shape of a Supabase Auth JWT. It never leaves this
// package — downstream code reads domain.AuthenticatedUser instead, which
// carries the app-resolved role rather than Supabase's own `role` claim
// (which is always "authenticated" for a signed-in user, not an RBAC role).
type supabaseClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	jwks     keyfunc.Keyfunc
	issuer   string
	audience string
	userSvc  userports.UserService
}

func NewAuthMiddleware(jwks keyfunc.Keyfunc, issuer, audience string, userSvc userports.UserService) *AuthMiddleware {
	return &AuthMiddleware{jwks: jwks, issuer: issuer, audience: audience, userSvc: userSvc}
}

func (m *AuthMiddleware) RequestAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken, err := extractBearer(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims := &supabaseClaims{}
		token, err := jwt.ParseWithClaims(rawToken, claims, m.jwks.Keyfunc,
			jwt.WithIssuer(m.issuer),
			jwt.WithAudience(m.audience),
			jwt.WithValidMethods([]string{"RS256", "ES256"}),
		)
		if err != nil || !token.Valid {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		authUserID := claims.Subject
		if authUserID == "" {
			http.Error(w, "invalid token: missing subject", http.StatusUnauthorized)
			return
		}

		authUser := &domain.AuthenticatedUser{
			AuthUserID: authUserID,
			Email:      claims.Email,
		}

		// A resolved local user isn't required to proceed: an authenticated
		// caller with no application user yet (first access) still reaches
		// the handler with UserID empty — only GET /me / POST /auth/profile
		// handle provisioning. Every other handler treats UserID=="" as
		// "not provisioned" via the existing owner/admin checks.
		if user, err := m.userSvc.GetByAuthUserID(r.Context(), authUserID); err == nil {
			authUser.UserID = user.ID
			authUser.Name = user.Name
			authUser.Email = user.Email
			authUser.Role = string(user.Role)
		}

		ctx := context.WithValue(r.Context(), claimsKey, authUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractBearer(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", &httpError{msg: "missing Authorization header"}
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", &httpError{msg: "Authorization header must be 'Bearer <token>'"}
	}
	return parts[1], nil
}

func GetClaims(ctx context.Context) *domain.AuthenticatedUser {
	claims, _ := ctx.Value(claimsKey).(*domain.AuthenticatedUser)
	return claims
}

// WithClaims returns a context carrying the given authenticated user.
// Intended for use in tests to simulate an authenticated request.
func WithClaims(ctx context.Context, claims *domain.AuthenticatedUser) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// IsAdmin reports whether the resolved app role is PLATFORM_ADMIN.
func IsAdmin(claims *domain.AuthenticatedUser) bool {
	return claims != nil && claims.Role == "PLATFORM_ADMIN"
}

type httpError struct{ msg string }

func (e *httpError) Error() string { return e.msg }
