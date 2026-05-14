package middleware

import (
	"context"
	"net/http"
	"strings"

	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
)

type contextKey string

const claimsKey contextKey = "claims"

// Extract reads Keycloak claims forwarded by KrakenD and stores them in the request context.
func Extract(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := authdomain.Claims{
			UserID: r.Header.Get("X-User-Id"),
			Email:  r.Header.Get("X-User-Email"),
			Name:   r.Header.Get("X-User-Name"),
			Roles:  parseRoles(r.Header.Get("X-User-Roles")),
		}
		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetClaims retrieves the Keycloak claims from the request context.
func GetClaims(ctx context.Context) (authdomain.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(authdomain.Claims)
	return claims, ok
}

// HasRole reports whether the claims in ctx include the given role.
func HasRole(ctx context.Context, role string) bool {
	claims, ok := GetClaims(ctx)
	if !ok {
		return false
	}
	for _, r := range claims.Roles {
		if strings.EqualFold(r, role) {
			return true
		}
	}
	return false
}

func parseRoles(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	roles := make([]string, 0, len(parts))
	for _, p := range parts {
		if r := strings.TrimSpace(p); r != "" {
			roles = append(roles, r)
		}
	}
	return roles
}
