package http

import (
	"encoding/json"
	"net/http"

	"github.com/diegoHDCz/ajudafio/internal/availability"
	"github.com/diegoHDCz/ajudafio/internal/availability/domain"
	"github.com/go-chi/chi/v5"
)

type AvailabilityHandler struct {
	s *availability.AvailabilityService
}

func NewAvailabilityHandler(s *availability.AvailabilityService) *AvailabilityHandler {
	return &AvailabilityHandler{s: s}
}

func NewAvailabilityRouter(h *AvailabilityHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/professional/{professionalID}", h.GetByProfessionalID)
	r.Post("/", h.Create)
	r.Patch("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

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

func (h *AvailabilityHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body createAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.ProfessionalID == "" || len(body.DayOfWeek) == 0 {
		http.Error(w, "professional_id and day_of_week are required", http.StatusBadRequest)
		return
	}
	a, err := h.s.Create(r.Context(), &domain.Availability{
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

func (h *AvailabilityHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
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

func (h *AvailabilityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.s.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

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
