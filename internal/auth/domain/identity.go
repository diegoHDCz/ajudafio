package domain

import "time"

type Identity struct {
	ID             string
	UserID         string
	Provider       string
	ProviderUserID string
	Email          string
	CreatedAt      time.Time
}

// GoogleClaims carries the subset of Google ID token claims GoogleLogin needs.
type GoogleClaims struct {
	Subject       string
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
}
