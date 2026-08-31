package http

import userdomain "github.com/diegoHDCz/ajudafio/internal/user/domain"

type meResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Role             string `json:"role"`
	OnboardingStatus string `json:"onboarding_status"`
}

func toMeResponse(u *userdomain.User) meResponse {
	return meResponse{
		ID:               u.ID,
		Name:             u.Name,
		Email:            u.Email,
		Role:             string(u.Role),
		OnboardingStatus: string(u.OnboardingStatus),
	}
}

type rolesResponse struct {
	Roles []string `json:"roles"`
}

type permissionsResponse struct {
	Permissions []string `json:"permissions"`
}

type completeProfileRequest struct {
	Role string `json:"role"`
}

type onboardingStatusRequest struct {
	Status string `json:"status"`
}

type onboardingStatusResponse struct {
	Status string `json:"status"`
}
