package http

import (
	"encoding/json"
	"net/http"

	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	"github.com/diegoHDCz/ajudafio/internal/availability"
	"github.com/diegoHDCz/ajudafio/internal/availability/domain"
	professionalports "github.com/diegoHDCz/ajudafio/internal/professional/ports"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AvailabilityHandler struct {
	s               *availability.AvailabilityService
	validator       *shared.Validator
	professionalSvc professionalports.ProfessionalService
}

func NewAvailabilityHandler(s *availability.AvailabilityService, validator *shared.Validator, professionalSvc professionalports.ProfessionalService) *AvailabilityHandler {
	return &AvailabilityHandler{s: s, validator: validator, professionalSvc: professionalSvc}
}

func NewAvailabilityRouter(h *AvailabilityHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/professional/{professionalID}", h.GetByProfessionalID)
	r.Post("/", h.Create)
	r.Patch("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

// @Summary      Buscar disponibilidades por profissional
// @Tags         availabilities
// @Produce      json
// @Param        professionalID  path      string  true  "Professional ID"
// @Success      200             {array}   availabilityResponse
// @Failure      500             {string}  string
// @Security     BearerAuth
// @Router       /availabilities/professional/{professionalID} [get]
func (h *AvailabilityHandler) GetByProfessionalID(w http.ResponseWriter, r *http.Request) {
	professionalID := chi.URLParam(r, "professionalID")
	list, err := h.s.GetByProfessionalID(r.Context(), professionalID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]availabilityResponse, len(list))
	for i, a := range list {
		resp[i] = toResponse(a)
	}
	respond(w, http.StatusOK, resp)
}

// @Summary      Criar disponibilidade
// @Tags         availabilities
// @Accept       json
// @Produce      json
// @Param        body  body      createAvailabilityRequest  true  "Dados da disponibilidade"
// @Success      201   {object}  availabilityResponse
// @Failure      400   {string}  string
// @Failure      422   {string}  string
// @Security     BearerAuth
// @Router       /availabilities [post]
func (h *AvailabilityHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body createAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.ProfessionalID == "" || body.DayOfWeek == "" || body.Shift == "" {
		http.Error(w, "professional_id, day_of_week and shift are required", http.StatusBadRequest)
		return
	}
	a, err := h.s.Create(r.Context(), &domain.Availability{
		ID:             uuid.New().String(),
		ProfessionalID: body.ProfessionalID,
		DayOfWeek:      body.DayOfWeek,
		Shift:          body.Shift,
		StartHour:      body.StartHour,
		EndHour:        body.EndHour,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	respond(w, http.StatusCreated, toResponse(a))
}

// @Summary      Atualizar disponibilidade
// @Tags         availabilities
// @Accept       json
// @Produce      json
// @Param        id    path      string                     true  "Availability ID"
// @Param        body  body      updateAvailabilityRequest  true  "Dados a atualizar"
// @Success      200   {object}  availabilityResponse
// @Failure      401   {string}  string
// @Failure      403   {string}  string
// @Failure      404   {string}  string
// @Failure      422   {string}  string
// @Security     BearerAuth
// @Router       /availabilities/{id} [patch]
func (h *AvailabilityHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	claims := authmiddleware.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !authmiddleware.IsAdmin(claims) {
		if err := h.checkOwnership(r, id); err != nil {
			http.Error(w, err.Error(), err.(httpErr).status)
			return
		}
	}

	var body updateAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	a, err := h.s.Update(r.Context(), &domain.Availability{
		ID:        id,
		DayOfWeek: body.DayOfWeek,
		Shift:     body.Shift,
		StartHour: body.StartHour,
		EndHour:   body.EndHour,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	respond(w, http.StatusOK, toResponse(a))
}

// @Summary      Remover disponibilidade
// @Tags         availabilities
// @Param        id  path  string  true  "Availability ID"
// @Success      204
// @Failure      401  {string}  string
// @Failure      403  {string}  string
// @Failure      404  {string}  string
// @Security     BearerAuth
// @Router       /availabilities/{id} [delete]
func (h *AvailabilityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	claims := authmiddleware.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.checkOwnership(r, id); err != nil {
		http.Error(w, err.Error(), err.(httpErr).status)
		return
	}

	if err := h.s.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// checkOwnership verifies the authenticated user owns the availability.
func (h *AvailabilityHandler) checkOwnership(r *http.Request, availabilityID string) error {
	claims := authmiddleware.GetClaims(r.Context())

	avail, err := h.s.GetByID(r.Context(), availabilityID)
	if err != nil {
		return httpErr{"availability not found", http.StatusNotFound}
	}

	p, err := h.professionalSvc.GetByID(r.Context(), avail.ProfessionalID)
	if err != nil {
		return httpErr{"professional not found", http.StatusNotFound}
	}

	if !h.validator.ValidateSameUserID(r.Context(), claims.Email, p.UserID) {
		return httpErr{"forbidden", http.StatusForbidden}
	}
	return nil
}

type httpErr struct {
	msg    string
	status int
}

func (e httpErr) Error() string { return e.msg }

func toResponse(a *domain.Availability) availabilityResponse {
	return availabilityResponse{
		ID:             a.ID,
		ProfessionalID: a.ProfessionalID,
		DayOfWeek:      a.DayOfWeek,
		Shift:          a.Shift,
		StartHour:      a.StartHour,
		EndHour:        a.EndHour,
	}
}

func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
