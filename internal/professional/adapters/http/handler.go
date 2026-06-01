package http

import (
	"encoding/json"
	"net/http"

	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	"github.com/diegoHDCz/ajudafio/internal/professional/ports"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	"github.com/go-chi/chi/v5"
)

type ProfessionalHandler struct {
	svc       ports.ProfessionalService
	validator *shared.Validator
}

func NewProfessionalHandler(svc ports.ProfessionalService, validator *shared.Validator) *ProfessionalHandler {
	return &ProfessionalHandler{svc: svc, validator: validator}
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

// @Summary      Buscar profissional por ID
// @Tags         professionals
// @Produce      json
// @Param        id   path      string  true  "Professional ID"
// @Success      200  {object}  professionalResponse
// @Failure      404  {string}  string
// @Router       /professionals/{id} [get]
func (h *ProfessionalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "professional not found", http.StatusNotFound)
		return
	}
	respond(w, http.StatusOK, toResponse(p))
}

// @Summary      Buscar profissional por User ID
// @Tags         professionals
// @Produce      json
// @Param        userID  path      string  true  "User ID"
// @Success      200     {object}  professionalResponse
// @Failure      404     {string}  string
// @Security     BearerAuth
// @Router       /professionals/user/{userID} [get]
func (h *ProfessionalHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	p, err := h.svc.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "professional not found", http.StatusNotFound)
		return
	}
	respond(w, http.StatusOK, toResponse(p))
}

// @Summary      Criar profissional
// @Tags         professionals
// @Accept       json
// @Produce      json
// @Param        body  body      createProfessionalRequest  true  "Dados do profissional"
// @Success      201   {object}  professionalResponse
// @Failure      400   {string}  string
// @Failure      422   {string}  string
// @Security     BearerAuth
// @Router       /professionals [post]
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

// @Summary      Atualizar profissional
// @Tags         professionals
// @Accept       json
// @Produce      json
// @Param        id    path      string                     true  "Professional ID"
// @Param        body  body      updateProfessionalRequest  true  "Dados a atualizar"
// @Success      200   {object}  professionalResponse
// @Failure      401   {string}  string
// @Failure      403   {string}  string
// @Failure      404   {string}  string
// @Failure      422   {string}  string
// @Security     BearerAuth
// @Router       /professionals/{id} [patch]
func (h *ProfessionalHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	claims := authmiddleware.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !authmiddleware.IsAdmin(claims) {
		p, err := h.svc.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, "professional not found", http.StatusNotFound)
			return
		}
		if !h.validator.ValidateSameUserID(r.Context(), claims.Email, p.UserID) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
	}

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

// @Summary      Remover profissional
// @Tags         professionals
// @Param        id  path  string  true  "Professional ID"
// @Success      204
// @Failure      401  {string}  string
// @Failure      403  {string}  string
// @Failure      404  {string}  string
// @Security     BearerAuth
// @Router       /professionals/{id} [delete]
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

	if !authmiddleware.IsAdmin(claims) && !h.validator.ValidateSameUserID(r.Context(), claims.Email, p.UserID) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Listar profissionais com filtros
// @Tags         professionals
// @Produce      json
// @Param        city         query     string    false  "Cidade"
// @Param        state        query     string    false  "Estado (UF)"
// @Param        day_of_week  query     []string  false  "Dias da semana"  collectionFormat(multi)
// @Param        shift        query     []string  false  "Turnos"          collectionFormat(multi)
// @Success      200  {array}   professionalResponse
// @Failure      500  {string}  string
// @Router       /professionals [get]
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
