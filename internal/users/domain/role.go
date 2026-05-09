package domain

import "fmt"

type Role string

const (
	RoleClient       Role = "CLIENT"
	RoleProfessional Role = "PROFESSIONAL"
	RoleAdmin        Role = "ADMIN"
)

func ParseRole(s string) (Role, error) {
	switch Role(s) {
	case RoleClient, RoleProfessional, RoleAdmin:
		return Role(s), nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidRole, s)
	}
}
