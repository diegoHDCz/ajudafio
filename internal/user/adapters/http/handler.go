package http

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
		http.Error(w, "User not found", http.StatusNotFound)
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

	// Como no banco Name é NOT NULL, validamos aqui ou deixamos o Service tratar
	if body.Name == "" || body.Email == "" {
		http.Error(w, "name and email are required", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Create(r.Context(), ports.CreateUserInput{
		Email: body.Email,
		Name:  body.Name,
		Phone: body.Phone,
		Role:  derefRole(body.Role, domain.RoleClient),
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
		ID:    id,
		Name:  body.Name,
		Email: body.Email,
		Phone: body.Phone,
		Role:  body.Role, // *string, o Service lida com a conversão e fallback para RoleClient
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



// ── Helpers ───────────────────────────────────────────────────────────────────

func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func derefString(s *string, fallback string) string {
	if s == nil {
		return fallback
	}
	return *s
}
func derefRole(role *domain.Role, fallback domain.Role) domain.Role {
	if role == nil {
		return fallback
	}
	return *role
}
