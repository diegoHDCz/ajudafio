package http

import "github.com/diegoHDCz/ajudafio/internal/user/domain"

type meResponse struct {
	Name  string      `json:"name"`
	Email string      `json:"email"`
	ID    string      `json:"id"`
	Role  domain.Role `json:"role"`
}

type createUserRequest struct {
	Email string       `json:"email"`
	Name  string       `json:"name"`
	Phone *string      `json:"phone"`
	Role  *domain.Role `json:"role"`
}

type updateUserRequest struct {
	Email *string      `json:"email"`
	Name  *string      `json:"name"`
	Phone *string      `json:"phone"`
	Role  *domain.Role `json:"role"`
}

type userResponse struct {
	ID        string      `json:"id"`
	Email     string      `json:"email"`
	Name      string      `json:"name"`
	Phone     *string     `json:"phone"`
	Role      domain.Role `json:"role"`
	AvatarURL *string     `json:"avatar_url,omitempty"`
	CreatedAt string      `json:"created_at"`
}

func toResponse(u *domain.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Phone:     u.Phone,
		Role:      u.Role,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
