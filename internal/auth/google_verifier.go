package auth

import (
	"context"
	"fmt"

	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	"google.golang.org/api/idtoken"
)

// GoogleTokenVerifier validates a Google ID token and extracts its claims.
// Abstracted behind an interface so service tests can substitute a fake
// instead of hitting Google's network/JWKS.
type GoogleTokenVerifier interface {
	Verify(ctx context.Context, idToken string) (*authdomain.GoogleClaims, error)
}

type googleTokenVerifier struct {
	clientID string
}

// NewGoogleTokenVerifier validates ID tokens against clientID using
// google.golang.org/api/idtoken, which handles Google's JWKS fetch/cache.
func NewGoogleTokenVerifier(clientID string) GoogleTokenVerifier {
	return &googleTokenVerifier{clientID: clientID}
}

func (v *googleTokenVerifier) Verify(ctx context.Context, idToken string) (*authdomain.GoogleClaims, error) {
	payload, err := idtoken.Validate(ctx, idToken, v.clientID)
	if err != nil {
		return nil, fmt.Errorf("google id token validation failed: %w", err)
	}

	email, _ := payload.Claims["email"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	return &authdomain.GoogleClaims{
		Subject:       payload.Subject,
		Email:         email,
		EmailVerified: emailVerified,
		Name:          name,
		Picture:       picture,
	}, nil
}
