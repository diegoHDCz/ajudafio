package domain

import "time"

type UserID string

type Role string

const (
	RoleClient       Role = "CLIENT"
	RoleProfessional Role = "PROFESSIONAL"
	RoleAdmin        Role = "ADMIN"
)

type User struct {
	ID                      UserID
	Email                   string
	Name                    *string
	EmailVerified           bool
	Image                   *string
	Telephone               *string
	TelephoneWhatsapp       bool
	SecondTelephone         *string
	SecondTelephoneWhatsapp bool
	Linkedin                *string
	Instagram               *string
	Facebook                *string
	IdentificationNumber    *string
	IdentificationType      *string
	Role                    Role
	CreatedAt               time.Time
	UpdatedAt               time.Time
}
