package http

import (
	"encoding/json"
	"net/http"

	"github.com/diegoHDCz/ajudafio/internal/auth"
	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *auth.AuthService
}

func NewHandler(svc *auth.AuthService) *Handler {
	return &Handler{svc: svc}
}

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(authmiddleware.Authenticate(h.svc))

	r.Get("/me", h.Me)
	r.Get("/me/accounts", h.MyAccounts)

	return r
}

// GET /auth/me
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := authmiddleware.GetUser(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	respond(w, http.StatusOK, toMeResponse(user))
}

// GET /auth/me/accounts
func (h *Handler) MyAccounts(w http.ResponseWriter, r *http.Request) {
	user, ok := authmiddleware.GetUser(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accounts, err := h.svc.GetAccountsByUser(r.Context(), string(user.ID))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := make([]accountResponse, len(accounts))
	for i, a := range accounts {
		resp[i] = toAccountResponse(a)
	}
	respond(w, http.StatusOK, resp)
}

func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
