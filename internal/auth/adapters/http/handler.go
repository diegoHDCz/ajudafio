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
	rp repository.KeycloakRepository
}

var (
	state string
)

func NewHandler(rp repository.KeycloakRepository) *Handler {
	state = "exemplo"
	return &Handler{rp: rp}
}

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/callback", h.Callback)

	return r
}

// GET /auth/callback
func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != state {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}
	config, err := h.rp.GetKeycloakConfig()
	oauth2Token, err := config.Exchange(r.Context(), r.URL.Query().Get("code"))

	if err != nil {
		http.Error(w, "Failed to exchange code", http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	log.Printf("rawIDToken: %s", rawIDToken)
	if !ok {
		http.Error(w, "Failed to get ID token", http.StatusInternalServerError)
		return
	}


	res := struct {
		OAuth2Token *oauth2.Token
		rawIDToken  string
	}{
		OAuth2Token: oauth2Token,
		rawIDToken:  rawIDToken,
	}

	respond(w, http.StatusOK, res)

}

func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
