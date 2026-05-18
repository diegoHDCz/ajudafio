package http

import (
	"encoding/json"
	"net/http"

	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	"github.com/diegoHDCz/ajudafio/internal/professional/ports"
	"github.com/go-chi/chi/v5"
)

type ProfessionalHandler struct {
	svc ports.ProfessionalService
}

func NewProfessionalHandler(svc ports.ProfessionalService) *ProfessionalHandler {
	return &ProfessionalHandler{svc: svc}
}

func NewRouter(handler *ProfessionalHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/", handler.FindWithFilters)
	r.Get("/user/{userID}", handler.GetByUserID)
	r.Get("/{id}", handler.GetByID)
	r.Post("/", handler.Create)
	r.Patch("/{id}", handler.Update)
	r.Delete("/{id}", handler.Delete)
	return r
}

func (h *ProfessionalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "professional not found", http.StatusNotFound)
		return
	}
	respond(w, http.StatusOK, toResponse(p))
}

func (h *ProfessionalHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	p, err := h.svc.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "professional not found", http.StatusNotFound)
		return
	}
	respond(w, http.StatusOK, toResponse(p))
}

func (h *ProfessionalHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body createProfessionalRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.UserID == "" || body.LicenseNumber == "" {
		http.Error(w, "user_id and license_number are required", http.StatusBadRequest)
		return
	}
	p, err := h.svc.Create(r.Context(), ports.CreateProfessionalInput{
		UserID:            body.UserID,
		LicenseNumber:     body.LicenseNumber,
		Category:          body.Category,
		YearsOfExperience: body.YearsOfExperience,
		Resume:            body.Resume,
		Metadata:          []byte(body.Metadata),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	respond(w, http.StatusCreated, toResponse(p))
}

func (h *ProfessionalHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body updateProfessionalRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	p, err := h.svc.Update(r.Context(), ports.UpdateProfessionalInput{
		ID:                id,
		LicenseNumber:     body.LicenseNumber,
		Category:          body.Category,
		YearsOfExperience: body.YearsOfExperience,
		Verified:          body.Verified,
		Resume:            body.Resume,
		Metadata:          []byte(body.Metadata),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	respond(w, http.StatusOK, toResponse(p))
}

func (h *ProfessionalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	claims := authmiddleware.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "professional not found", http.StatusNotFound)
		return
	}

	isAdmin := false
	for _, role := range claims.RealmAccess.Roles {
		if role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin && claims.Sub != p.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProfessionalHandler) FindWithFilters(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filters := ports.ProfessionalFilters{}
	if city := q.Get("city"); city != "" {
		filters.City = &city
	}
	if state := q.Get("state"); state != "" {
		filters.State = &state
	}
	if days := q["day_of_week"]; len(days) > 0 {
		filters.DayOfWeek = days
	}
	if shifts := q["shift"]; len(shifts) > 0 {
		filters.Shift = shifts
	}
	list, err := h.svc.FindWithFilters(r.Context(), filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]professionalResponse, len(list))
	for i, p := range list {
		resp[i] = toResponse(p)
	}
	respond(w, http.StatusOK, resp)
}

func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
