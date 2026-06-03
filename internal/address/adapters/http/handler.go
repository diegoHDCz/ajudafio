package http

import (
	"encoding/json"
	"net/http"

	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	"github.com/diegoHDCz/ajudafio/internal/address/ports"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	"github.com/go-chi/chi/v5"
)

type AddressHandler struct {
	svc       ports.AddressService
	validator *shared.Validator
}

func NewAddressHandler(svc ports.AddressService, validator *shared.Validator) *AddressHandler {
	return &AddressHandler{svc: svc, validator: validator}
}

func NewRouter(handler *AddressHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", handler.GetByID)
	r.Get("/user/{userID}", handler.GetByUserID)
	r.Post("/", handler.Create)
	r.Patch("/{id}", handler.Update)
	r.Delete("/{id}", handler.Delete)
	return r
}

// @Summary      Buscar endereço por ID
// @Tags         addresses
// @Produce      json
// @Param        id   path      string  true  "Address ID"
// @Success      200  {object}  addressResponse
// @Failure      404  {string}  string
// @Security     BearerAuth
// @Router       /addresses/{id} [get]
func (h *AddressHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	address, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "address not found", http.StatusNotFound)
		return
	}
	respond(w, http.StatusOK, toResponse(address))
}

// @Summary      Buscar endereços por User ID
// @Tags         addresses
// @Produce      json
// @Param        userID  path      string  true  "User ID"
// @Success      200     {array}   addressResponse
// @Failure      500     {string}  string
// @Security     BearerAuth
// @Router       /addresses/user/{userID} [get]
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

// @Summary      Criar endereço
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        body  body      createAddressRequest  true  "Dados do endereço"
// @Success      201   {object}  addressResponse
// @Failure      400   {string}  string
// @Failure      422   {string}  string
// @Security     BearerAuth
// @Router       /addresses [post]
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

// @Summary      Atualizar endereço
// @Tags         addresses
// @Accept       json
// @Produce      json
// @Param        id    path      string                true  "Address ID"
// @Param        body  body      updateAddressRequest  true  "Dados a atualizar"
// @Success      200   {object}  addressResponse
// @Failure      401   {string}  string
// @Failure      403   {string}  string
// @Failure      404   {string}  string
// @Failure      422   {string}  string
// @Security     BearerAuth
// @Router       /addresses/{id} [patch]
func (h *AddressHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	claims := authmiddleware.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !authmiddleware.IsAdmin(claims) {
		a, err := h.svc.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, "address not found", http.StatusNotFound)
			return
		}
		if !h.validator.ValidateSameUserID(r.Context(), claims.Email, a.UserID) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
	}

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

// @Summary      Remover endereço
// @Tags         addresses
// @Param        id  path  string  true  "Address ID"
// @Success      204
// @Failure      401  {string}  string
// @Failure      403  {string}  string
// @Failure      404  {string}  string
// @Security     BearerAuth
// @Router       /addresses/{id} [delete]
func (h *AddressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	claims := authmiddleware.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !authmiddleware.IsAdmin(claims) {
		a, err := h.svc.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, "address not found", http.StatusNotFound)
			return
		}
		if !h.validator.ValidateSameUserID(r.Context(), claims.Email, a.UserID) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
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
