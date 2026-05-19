package http

import (
	"context"
	"encoding/json"

	"net/http"

	repository "github.com/diegoHDCz/ajudafio/internal/auth/adapters/keycloak"
	tokens "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	roles "github.com/diegoHDCz/ajudafio/internal/user/domain"
	"github.com/diegoHDCz/ajudafio/internal/user/ports"
	"github.com/go-chi/chi/v5"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
)

type Handler struct {
	rp     repository.KeycloakRepository
	config *oauth2.Config
	us     ports.UserService
}

var (
	state string
)

func NewHandler(rp repository.KeycloakRepository, config *oauth2.Config, us ports.UserService) *Handler {
	state = "exemplo"
	return &Handler{rp: rp, config: config, us: us}
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

	res := struct {
		OAuth2Token *oauth2.Token
		rawIDToken  string
	}{
		OAuth2Token: oauth2Token,
		rawIDToken:  rawIDToken,
	}
	h.syncDBuser(r.Context(), oauth2Token)
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

func (h *Handler) syncDBuser(ctx context.Context, token *oauth2.Token) error {
	jwks, err := tokens.InitJWKS(ctx)
	if err != nil {
		log.Printf("Failed to initialize JWKS: %v", err)
		return err
	}
	claims, err := tokens.ParseToken(token.AccessToken, jwks)

	user, err := h.us.GetByEmail(ctx, claims.Email)

	if err != nil || user == nil {
		_, err = h.us.Create(ctx, ports.CreateUserInput{
			ID:    claims.Sub,
			Name:  claims.Name,
			Email: claims.Email,
			Role:  roles.RoleClient,
			Phone: nil,
		})
		if err != nil {
			log.Printf("Failed to create user: %v", err)
			return err
		}
		log.Printf("User created: %v", claims.Email)
	}

	return nil
}
