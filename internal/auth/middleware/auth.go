package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/diegoHDCz/ajudafio/internal/auth"
	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	userdomain "github.com/diegoHDCz/ajudafio/internal/user/domain"
)

type contextKey string

const (
	claimsKey contextKey = "claims"
	userKey   contextKey = "user"
)

// Extract reads Keycloak claims forwarded by KrakenD and stores them in the request context.
// It does not enforce authentication — use Authenticate for that.
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

// Authenticate requires valid Keycloak claims and hydrates the local user from the
// database, creating the user row and account link on first login.
// Returns 401 if claims are missing and 500 on a database error.
func Authenticate(svc *auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := authdomain.Claims{
				UserID: r.Header.Get("X-User-Id"),
				Email:  r.Header.Get("X-User-Email"),
				Name:   r.Header.Get("X-User-Name"),
				Roles:  parseRoles(r.Header.Get("X-User-Roles")),
			}

			if claims.UserID == "" || claims.Email == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := svc.FindOrCreateUser(r.Context(), claims)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			ctx = context.WithValue(ctx, userKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole returns a middleware that allows the request only if the authenticated
// user has the given role. Must be chained after Authenticate.
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !HasRole(r.Context(), role) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r.WithContext(r.Context()))
		})
	}
}

// GetClaims retrieves the Keycloak claims from the request context.
func GetClaims(ctx context.Context) (authdomain.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(authdomain.Claims)
	return claims, ok
}

// GetUser retrieves the local user hydrated by Authenticate from the request context.
func GetUser(ctx context.Context) (*userdomain.User, bool) {
	user, ok := ctx.Value(userKey).(*userdomain.User)
	return user, ok
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
