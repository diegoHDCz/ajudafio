package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/diegoHDCz/ajudafio/internal/auth"
	"github.com/diegoHDCz/ajudafio/internal/auth/ports"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc ports.AuthService
}

func NewHandler(svc ports.AuthService) *Handler {
	return &Handler{svc: svc}
}

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	r.Post("/logout", h.Logout)
	return r
}

// @Summary      Registrar novo usuário
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      registerRequest  true  "Dados de cadastro"
// @Success      201   {object}  tokenResponse
// @Failure      400   {string}  string
// @Failure      409   {string}  string
// @Router       /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var body registerRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Name == "" || body.Email == "" || body.Password == "" {
		http.Error(w, "name, email and password are required", http.StatusBadRequest)
		return
	}

	email, err := shared.NormalizeEmail(body.Email)
	if err != nil {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	pair, err := h.svc.Register(r.Context(), ports.RegisterInput{
		Name:     body.Name,
		Email:    email,
		Phone:    body.Phone,
		Password: body.Password,
	})
	if err != nil {
		if errors.Is(err, auth.ErrEmailAlreadyInUse) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusCreated, toTokenResponse(pair))
}

// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      loginRequest  true  "Credenciais"
// @Success      200   {object}  tokenResponse
// @Failure      400   {string}  string
// @Failure      401   {string}  string
// @Router       /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body loginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Email == "" || body.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	email, err := shared.NormalizeEmail(body.Email)
	if err != nil {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	pair, err := h.svc.Login(r.Context(), email, body.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	respond(w, http.StatusOK, toTokenResponse(pair))
}

// @Summary      Renovar access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      refreshRequest  true  "Refresh token"
// @Success      200   {object}  tokenResponse
// @Failure      400   {string}  string
// @Failure      401   {string}  string
// @Router       /auth/refresh [post]
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.RefreshToken == "" {
		http.Error(w, "refresh_token is required", http.StatusBadRequest)
		return
	}

	pair, err := h.svc.Refresh(r.Context(), body.RefreshToken)
	if err != nil {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	respond(w, http.StatusOK, toTokenResponse(pair))
}

// @Summary      Logout
// @Tags         auth
// @Accept       json
// @Param        body  body  logoutRequest  true  "Refresh token"
// @Success      204
// @Failure      400  {string}  string
// @Router       /auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var body logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.RefreshToken == "" {
		http.Error(w, "refresh_token is required", http.StatusBadRequest)
		return
	}

	if err := h.svc.Logout(r.Context(), body.RefreshToken); err != nil {
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
