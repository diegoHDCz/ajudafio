package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	oidc "github.com/coreos/go-oidc"
	"github.com/diegoHDCz/ajudafio/internal/auth/domain"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const claimsKey contextKey = "claims"

type AuthMiddleware struct {
	verifier *oidc.IDTokenVerifier
}

func NewAuthMiddleware(ctx context.Context, issuerURL, clientID string) (*AuthMiddleware, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, err
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	return &AuthMiddleware{verifier: verifier}, nil
}

// RequestAuth validates the Bearer token and stores claims in the request context.
func (m *AuthMiddleware) RequestAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwks, err := InitJWKS(r.Context())
		if err != nil {
			http.Error(w, "failed to initialize JWKS", http.StatusInternalServerError)
			return
		}

		rawToken, err := extractBearer(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := ParseToken(rawToken, jwks)

		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func InitJWKS(ctx context.Context) (keyfunc.Keyfunc, error) {
	jwksURL := "http://localhost:8180/realms/ajudafio/protocol/openid-connect/certs"

	var err error
	jwks, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		log.Printf("failed to create JWKS from URL: %v", err)
		return nil, err
	}
	return jwks, nil
}

func ParseToken(tokenString string, jwks keyfunc.Keyfunc) (*domain.JWTClaims, error) {

	claims := &domain.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, jwks.Keyfunc)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token inválido: %w", err)
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

type httpError struct{ msg string }

func (e *httpError) Error() string { return e.msg }
