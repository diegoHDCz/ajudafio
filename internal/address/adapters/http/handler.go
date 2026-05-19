package http

import (
	"encoding/json"
	"net/http"

	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	"github.com/diegoHDCz/ajudafio/internal/address/ports"
	"github.com/go-chi/chi/v5"
)

type AddressHandler struct {
	svc ports.AddressService
}

func NewAddressHandler(svc ports.AddressService) *AddressHandler {
	return &AddressHandler{svc: svc}
}

func NewRouter(handler *AddressHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", handler.GetByID)
	r.Get("/user/{userID}", handler.GetByUserID)
	r.Get("/contract/{contractID}", handler.GetByContractID)
	r.Post("/", handler.Create)
	r.Patch("/{id}", handler.Update)
	r.Delete("/{id}", handler.Delete)
	return r
}

func (h *AddressHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	address, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "address not found", http.StatusNotFound)
		return
	}
	respond(w, http.StatusOK, toResponse(address))
}

func (h *AddressHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	addresses, err := h.svc.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]addressResponse, len(addresses))
	for i, a := range addresses {
		resp[i] = toResponse(a)
	}
	respond(w, http.StatusOK, resp)
}

func (h *AddressHandler) GetByContractID(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "contractID")
	addresses, err := h.svc.GetByContractID(r.Context(), contractID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]addressResponse, len(addresses))
	for i, a := range addresses {
		resp[i] = toResponse(a)
	}
	respond(w, http.StatusOK, resp)
}

func (h *AddressHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body createAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.UserID == "" || body.ZipCode == "" || body.AddressLine == "" || body.Number == "" || body.District == "" || body.City == "" || body.State == "" {
		http.Error(w, "user_id, zip_code, address_line, number, district, city and state are required", http.StatusBadRequest)
		return
	}
	address, err := h.svc.Create(r.Context(), ports.CreateAddressInput{
		UserID:      body.UserID,
		ContractID:  body.ContractID,
		ZipCode:     body.ZipCode,
		AddressLine: body.AddressLine,
		Number:      body.Number,
		Complement:  body.Complement,
		District:    body.District,
		City:        body.City,
		State:       body.State,
		Reference:   body.Reference,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	respond(w, http.StatusCreated, toResponse(address))
}

func (h *AddressHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body updateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	address, err := h.svc.Update(r.Context(), ports.UpdateAddressInput{
		ID:          id,
		ZipCode:     body.ZipCode,
		AddressLine: body.AddressLine,
		Number:      body.Number,
		Complement:  body.Complement,
		District:    body.District,
		City:        body.City,
		State:       body.State,
		Reference:   body.Reference,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	respond(w, http.StatusOK, toResponse(address))
}

func (h *AddressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	claims := authmiddleware.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	isAdmin := false
	for _, role := range claims.RealmAccess.Roles {
		if role == "admin" {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
