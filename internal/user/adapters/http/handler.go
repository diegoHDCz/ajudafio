package userhttp

import (
	"encoding/json"
	"net/http"

	"github.com/diegoHDCz/ajudafio/internal/user/domain"
	"github.com/diegoHDCz/ajudafio/internal/user/ports"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc ports.UserService
}

func NewHandler(svc ports.UserService) *Handler {
	return &Handler{svc: svc}
}

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", h.GetByID)
	r.Post("/", h.Create)
	r.Patch("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

// GET /users/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(chi.URLParam(r, "id"))

	user, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	respond(w, http.StatusOK, toResponse(user))
}

// POST /users
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Create(r.Context(), ports.CreateUserInput{
		Email:                   body.Email,
		Name:                    body.Name,
		Telephone:               body.Telephone,
		TelephoneWhatsapp:       derefBool(body.TelephoneWhatsapp),
		SecondTelephone:         body.SecondTelephone,
		SecondTelephoneWhatsapp: derefBool(body.SecondTelephoneWhatsapp),
		Linkedin:                body.Linkedin,
		Instagram:               body.Instagram,
		Facebook:                body.Facebook,
		IdentificationNumber:    body.IdentificationNumber,
		IdentificationType:      body.IdentificationType,
		Role:                    domain.Role(derefString(body.Role, string(domain.RoleClient))),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusCreated, toResponse(user))
}

// PATCH /users/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(chi.URLParam(r, "id"))

	var body updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Update(r.Context(), ports.UpdateUserInput{
		ID:                      id,
		Name:                    body.Name,
		Telephone:               body.Telephone,
		TelephoneWhatsapp:       body.TelephoneWhatsapp,
		SecondTelephone:         body.SecondTelephone,
		SecondTelephoneWhatsapp: body.SecondTelephoneWhatsapp,
		Linkedin:                body.Linkedin,
		Instagram:               body.Instagram,
		Facebook:                body.Facebook,
		IdentificationNumber:    body.IdentificationNumber,
		IdentificationType:      body.IdentificationType,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusOK, toResponse(user))
}

// DELETE /users/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(chi.URLParam(r, "id"))

	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ── DTOs ─────────────────────────────────────────────────────────────────────

type createUserRequest struct {
	Email                   string  `json:"email"`
	Name                    *string `json:"name"`
	Telephone               *string `json:"telephone"`
	TelephoneWhatsapp       *bool   `json:"telephone_whatsapp"`
	SecondTelephone         *string `json:"second_telephone"`
	SecondTelephoneWhatsapp *bool   `json:"second_telephone_whatsapp"`
	Linkedin                *string `json:"linkedin"`
	Instagram               *string `json:"instagram"`
	Facebook                *string `json:"facebook"`
	IdentificationNumber    *string `json:"identification_number"`
	IdentificationType      *string `json:"identification_type"`
	Role                    *string `json:"role"`
}

type updateUserRequest struct {
	Name                    *string `json:"name"`
	Telephone               *string `json:"telephone"`
	TelephoneWhatsapp       *bool   `json:"telephone_whatsapp"`
	SecondTelephone         *string `json:"second_telephone"`
	SecondTelephoneWhatsapp *bool   `json:"second_telephone_whatsapp"`
	Linkedin                *string `json:"linkedin"`
	Instagram               *string `json:"instagram"`
	Facebook                *string `json:"facebook"`
	IdentificationNumber    *string `json:"identification_number"`
	IdentificationType      *string `json:"identification_type"`
}

type userResponse struct {
	ID                      string  `json:"id"`
	Email                   string  `json:"email"`
	Name                    *string `json:"name"`
	EmailVerified           bool    `json:"email_verified"`
	Telephone               *string `json:"telephone"`
	TelephoneWhatsapp       bool    `json:"telephone_whatsapp"`
	SecondTelephone         *string `json:"second_telephone"`
	SecondTelephoneWhatsapp bool    `json:"second_telephone_whatsapp"`
	Linkedin                *string `json:"linkedin"`
	Instagram               *string `json:"instagram"`
	Facebook                *string `json:"facebook"`
	IdentificationNumber    *string `json:"identification_number"`
	IdentificationType      *string `json:"identification_type"`
	Role                    string  `json:"role"`
}

func toResponse(u *domain.User) userResponse {
	return userResponse{
		ID:                      string(u.ID),
		Email:                   u.Email,
		Name:                    u.Name,
		EmailVerified:           u.EmailVerified,
		Telephone:               u.Telephone,
		TelephoneWhatsapp:       u.TelephoneWhatsapp,
		SecondTelephone:         u.SecondTelephone,
		SecondTelephoneWhatsapp: u.SecondTelephoneWhatsapp,
		Linkedin:                u.Linkedin,
		Instagram:               u.Instagram,
		Facebook:                u.Facebook,
		IdentificationNumber:    u.IdentificationNumber,
		IdentificationType:      u.IdentificationType,
		Role:                    string(u.Role),
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func derefString(s *string, fallback string) string {
	if s == nil {
		return fallback
	}
	return *s
}
