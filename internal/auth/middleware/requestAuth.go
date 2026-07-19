package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	auth "github.com/diegoHDCz/ajudafio/internal/auth"
	"github.com/diegoHDCz/ajudafio/internal/auth/domain"
)

type contextKey string

const claimsKey contextKey = "claims"

const minSecretLen = 32

type AuthMiddleware struct {
	secret []byte
}

func NewAuthMiddleware(secret []byte) (*AuthMiddleware, error) {
	if len(secret) < minSecretLen {
		return nil, fmt.Errorf("JWT secret must be at least %d bytes", minSecretLen)
	}
	return &AuthMiddleware{secret: secret}, nil
}

func (m *AuthMiddleware) RequestAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken, err := extractBearer(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateAccessToken(m.secret, rawToken)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
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

func GetClaims(ctx context.Context) *domain.JWTClaims {
	claims, _ := ctx.Value(claimsKey).(*domain.JWTClaims)
	return claims
}

// WithClaims returns a context carrying the given JWT claims.
// Intended for use in tests to simulate an authenticated request.
func WithClaims(ctx context.Context, claims *domain.JWTClaims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// IsAdmin reports whether the claims include the "ADMIN" role.
func IsAdmin(claims *domain.JWTClaims) bool {
	return claims.Role == "ADMIN"
}

type httpError struct{ msg string }

func (e *httpError) Error() string { return e.msg }
