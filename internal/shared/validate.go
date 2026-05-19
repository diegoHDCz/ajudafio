package shared

import (
	"context"

	"github.com/diegoHDCz/ajudafio/internal/user/ports"
)

type Validator struct {
	svc ports.UserService
}

func NewValidator(svc ports.UserService) *Validator {
	return &Validator{svc: svc}
}

func (v *Validator) ValidateSameUserID(context context.Context, email string, userID2 string) bool {
	user, err := v.svc.GetByEmail(context, email)
	if err != nil {
		return false
	}
	return string(user.ID) == userID2
}
