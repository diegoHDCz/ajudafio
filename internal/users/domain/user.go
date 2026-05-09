package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/seuuser/healthcontracts/internal/shared/valueobjects"
)

var (
	ErrEmailAlreadyInUse  = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidRole        = errors.New("invalid role")
)

type User struct {
	id                   valueobjects.EntityID
	email                string
	name                 string
	emailVerified        bool
	role                 Role
	telephone            *string
	identificationNumber *string
	createdAt            time.Time
	updatedAt            time.Time
}

func NewUser(email, name string) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	now := time.Now().UTC()
	return &User{
		id:            valueobjects.NewEntityID(),
		email:         strings.ToLower(email),
		name:          name,
		emailVerified: false,
		role:          RoleClient,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

// Reconstituição a partir do banco — sem validações de criação.
func Reconstitute(
	id valueobjects.EntityID,
	email, name string,
	emailVerified bool,
	role Role,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id: id, email: email, name: name,
		emailVerified: emailVerified, role: role,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}

func (u *User) ID() valueobjects.EntityID { return u.id }
func (u *User) Email() string             { return u.email }
func (u *User) Name() string              { return u.name }
func (u *User) Role() Role                { return u.role }
func (u *User) EmailVerified() bool       { return u.emailVerified }
