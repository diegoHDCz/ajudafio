package http

import (
	"encoding/json"

	"net/http"

	repository "github.com/diegoHDCz/ajudafio/internal/auth/adapters/keycloak"
	"github.com/go-chi/chi/v5"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
)

type Handler struct {
	rp     repository.KeycloakRepository
	config *oauth2.Config
}

var (
	state string
)

func NewHandler(rp repository.KeycloakRepository, config *oauth2.Config) *Handler {
	state = "exemplo"
	return &Handler{rp: rp, config: config}
}

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Get("/login", h.Login)
	r.Get("/callback", h.Callback)

	return r
}

// GET /auth/callback
func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != state {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	oauth2Token, err := h.config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {

		http.Error(w, "Failed to validate token", http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "Failed to get ID token", http.StatusInternalServerError)
		return
	}

	log.Printf("teste raw %v", rawIDToken)

	res := struct {
		OAuth2Token *oauth2.Token
		rawIDToken  string
	}{
		OAuth2Token: oauth2Token,
		rawIDToken:  rawIDToken,
	}

	respond(w, http.StatusOK, res)

}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, h.config.AuthCodeURL("exemplo"), http.StatusFound)
}

func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
