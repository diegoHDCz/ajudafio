package domain

// UserAuth carries the credential data needed to authenticate a user.
// It is internal to the auth module and must never be serialized in an HTTP response.
type UserAuth struct {
	ID           string
	Email        string
	Name         string
	Role         string
	PasswordHash string
}
