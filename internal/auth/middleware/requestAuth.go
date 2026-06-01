package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/diegoHDCz/ajudafio/internal/auth/domain"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const claimsKey contextKey = "claims"

type AuthMiddleware struct {
	jwks keyfunc.Keyfunc
}

func NewAuthMiddleware(ctx context.Context, jwksURL string) (*AuthMiddleware, error) {
	jwks, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JWKS from %s: %w", jwksURL, err)
	}
	return &AuthMiddleware{jwks: jwks}, nil
}

func (m *AuthMiddleware) RequestAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken, err := extractBearer(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := parseToken(rawToken, m.jwks)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseToken(tokenString string, jwks keyfunc.Keyfunc) (*domain.JWTClaims, error) {
	claims := &domain.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, jwks.Keyfunc)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	return claims, nil
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

// IsAdmin reports whether the claims include the "admin" role.
func IsAdmin(claims *domain.JWTClaims) bool {
	return claims.Role == "admin"
}

type httpError struct{ msg string }

func (e *httpError) Error() string { return e.msg }
