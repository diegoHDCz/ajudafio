package domain

// Claims holds the Keycloak user data forwarded by KrakenD as request headers.
type Claims struct {
	UserID string
	Email  string
	Name   string
	Roles  []string
}
