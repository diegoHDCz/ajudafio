package http

import (
	"encoding/json"
	"net/http"

	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	"github.com/diegoHDCz/ajudafio/internal/user/domain"
	"github.com/diegoHDCz/ajudafio/internal/user/ports"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc       ports.UserService
	validator *shared.Validator
}

func NewHandler(svc ports.UserService, validator *shared.Validator) *Handler {
	return &Handler{svc: svc, validator: validator}
}

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Get("/me", h.Me)
	r.Get("/{id}", h.GetByID)
	r.Patch("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

// @Summary      Dados do usuário autenticado
// @Tags         users
// @Produce      json
// @Success      200  {object}  meResponse
// @Failure      401  {string}  string
// @Security     BearerAuth
// @Router       /users/me [get]
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims := authmiddleware.GetClaims(r.Context())

	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	respond(w, http.StatusOK, meResponse{
		Name:  claims.Name,
		Email: claims.Email,
	})
}

// @Summary      Buscar usuário por ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  userResponse
// @Failure      404  {string}  string
// @Security     BearerAuth
// @Router       /users/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(chi.URLParam(r, "id"))

	user, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	respond(w, http.StatusOK, toResponse(user))
}

// @Summary      Criar usuário
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      createUserRequest  true  "Dados do usuário"
// @Success      201   {object}  userResponse
// @Failure      400   {string}  string
// @Failure      500   {string}  string
// @Security     BearerAuth
// @Router       /users [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

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

// @Summary      Atualizar usuário
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string             true  "User ID"
// @Param        body  body      updateUserRequest  true  "Dados a atualizar"
// @Success      200   {object}  userResponse
// @Failure      400   {string}  string
// @Failure      401   {string}  string
// @Failure      403   {string}  string
// @Failure      500   {string}  string
// @Security     BearerAuth
// @Router       /users/{id} [patch]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(chi.URLParam(r, "id"))

	claims := authmiddleware.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !authmiddleware.IsAdmin(claims) && !h.validator.ValidateSameUserID(r.Context(), claims.Email, string(id)) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

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
		Role:  body.Role,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusOK, toResponse(user))
}

// @Summary      Remover usuário
// @Tags         users
// @Param        id  path  string  true  "User ID"
// @Success      204
// @Failure      401  {string}  string
// @Failure      403  {string}  string
// @Failure      500  {string}  string
// @Security     BearerAuth
// @Router       /users/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(chi.URLParam(r, "id"))

	claims := authmiddleware.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !authmiddleware.IsAdmin(claims) && !h.validator.ValidateSameUserID(r.Context(), claims.Email, string(id)) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

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

func derefRole(role *domain.Role, fallback domain.Role) domain.Role {
	if role == nil {
		return fallback
	}
	return *role
}
