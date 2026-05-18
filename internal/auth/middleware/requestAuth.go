package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	oidc "github.com/coreos/go-oidc"
)

type contextKey string

const claimsKey contextKey = "claims"

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type ClientAccess struct {
	Roles []string `json:"roles"`
}

type Claims struct {
	Subject        string                  `json:"sub"`
	Email          string                  `json:"email"`
	Name           string                  `json:"name"`
	PreferredName  string                  `json:"preferred_username"`
	RealmAccess    RealmAccess             `json:"realm_access"`
	ResourceAccess map[string]ClientAccess `json:"resource_access"`
}

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
		rawToken, err := extractBearer(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		idToken, err := m.verifier.Verify(r.Context(), rawToken)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		var claims Claims
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "failed to parse token claims", http.StatusInternalServerError)
			return
		}
		log.Printf("tokens %v", idToken)
		log.Printf("claims %v", claims)
		ctx := context.WithValue(r.Context(), claimsKey, &claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole returns a middleware that enforces a realm-level role.
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !HasRealmRole(r.Context(), role) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireClientRole returns a middleware that enforces a client-level role.
func RequireClientRole(clientID, role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !HasClientRole(r.Context(), clientID, role) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// GetClaims retrieves the parsed claims from the request context.
func GetClaims(ctx context.Context) *Claims {
	claims, _ := ctx.Value(claimsKey).(*Claims)
	return claims
}

// HasRealmRole reports whether the caller has the given Keycloak realm role.
func HasRealmRole(ctx context.Context, role string) bool {
	claims := GetClaims(ctx)
	if claims == nil {
		return false
	}
	for _, r := range claims.RealmAccess.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasClientRole reports whether the caller has the given role for a specific client.
func HasClientRole(ctx context.Context, clientID, role string) bool {
	claims := GetClaims(ctx)
	if claims == nil {
		return false
	}
	access, ok := claims.ResourceAccess[clientID]
	if !ok {
		return false
	}
	for _, r := range access.Roles {
		if r == role {
			return true
		}
	}
	return false
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
	log.Printf("parts %v", len(parts))
	return parts[1], nil
}

type httpError struct{ msg string }

func (e *httpError) Error() string { return e.msg }
