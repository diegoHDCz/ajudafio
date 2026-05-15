package http

import (
	"github.com/diegoHDCz/ajudafio/internal/auth/ports"
	userdomain "github.com/diegoHDCz/ajudafio/internal/user/domain"
)

type meResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Phone     *string `json:"phone,omitempty"`
	Role      string  `json:"role"`
	CreatedAt string  `json:"created_at"`
}

type accountResponse struct {
	ID         string `json:"id"`
	ProviderID string `json:"provider_id"`
	AccountID  string `json:"account_id"`
	CreatedAt  string `json:"created_at"`
}

func toMeResponse(u *userdomain.User) meResponse {
	return meResponse{
		ID:        string(u.ID),
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toAccountResponse(a ports.Account) accountResponse {
	return accountResponse{
		ID:         a.ID,
		ProviderID: a.ProviderID,
		AccountID:  a.AccountID,
		CreatedAt:  a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
