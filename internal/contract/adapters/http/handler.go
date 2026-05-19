package http

import (
	"encoding/json"
	"net/http"
	"time"

	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	"github.com/diegoHDCz/ajudafio/internal/contract/domain"
	"github.com/diegoHDCz/ajudafio/internal/contract/ports"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	"github.com/go-chi/chi/v5"
)

type ContractHandler struct {
	svc       ports.ContractService
	validator *shared.Validator
}

func NewContractHandler(svc ports.ContractService, validator *shared.Validator) *ContractHandler {
	return &ContractHandler{svc: svc, validator: validator}
}

func NewRouter(h *ContractHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.GetAll)
	r.Get("/{id}", h.GetByID)
	r.Get("/user/{userID}", h.GetByUserID)
	r.Get("/professional/{professionalID}", h.GetByProfessionalID)
	r.Get("/status/{status}", h.GetByStatus)
	r.Get("/category/{category}", h.GetByProfessionalCategory)
	r.Post("/", h.Create)
	r.Patch("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

func (h *ContractHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	contracts, err := h.svc.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]contractResponse, len(contracts))
	for i, c := range contracts {
		resp[i] = toResponse(c)
	}
	respond(w, http.StatusOK, resp)
}

func (h *ContractHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	c, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "contract not found", http.StatusNotFound)
		return
	}
	respond(w, http.StatusOK, toResponse(c))
}

func (h *ContractHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	status := r.URL.Query().Get("status")

	var (
		contracts []*domain.Contract
		err       error
	)
	if status != "" {
		contracts, err = h.svc.GetByUserIDAndStatus(r.Context(), userID, status)
	} else {
		contracts, err = h.svc.GetByUserID(r.Context(), userID)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]contractResponse, len(contracts))
	for i, c := range contracts {
		resp[i] = toResponse(c)
	}
	respond(w, http.StatusOK, resp)
}

func (h *ContractHandler) GetByProfessionalID(w http.ResponseWriter, r *http.Request) {
	professionalID := chi.URLParam(r, "professionalID")
	status := r.URL.Query().Get("status")

	var (
		contracts []*domain.Contract
		err       error
	)
	if status != "" {
		contracts, err = h.svc.GetByProfessionalIDAndStatus(r.Context(), professionalID, status)
	} else {
		contracts, err = h.svc.GetByProfessionalID(r.Context(), professionalID)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]contractResponse, len(contracts))
	for i, c := range contracts {
		resp[i] = toResponse(c)
	}
	respond(w, http.StatusOK, resp)
}

func (h *ContractHandler) GetByStatus(w http.ResponseWriter, r *http.Request) {
	status := chi.URLParam(r, "status")
	contracts, err := h.svc.GetByStatus(r.Context(), status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]contractResponse, len(contracts))
	for i, c := range contracts {
		resp[i] = toResponse(c)
	}
	respond(w, http.StatusOK, resp)
}

func (h *ContractHandler) GetByProfessionalCategory(w http.ResponseWriter, r *http.Request) {
	category := chi.URLParam(r, "category")
	contracts, err := h.svc.GetByProfessionalCategory(r.Context(), category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]contractResponse, len(contracts))
	for i, c := range contracts {
		resp[i] = toResponse(c)
	}
	respond(w, http.StatusOK, resp)
}

func (h *ContractHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body createContractRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.ClientID == "" || body.ProfessionalID == "" {
		http.Error(w, "client_id and professional_id are required", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse("15:04", body.StartTime)
	if err != nil {
		http.Error(w, "start_time must be in HH:MM format", http.StatusBadRequest)
		return
	}

	weekDays := shared.SliceToDayOfWeek(body.WeekDays)

	contract, err := h.svc.Create(r.Context(), ports.CreateContractInput{
		ClientID:       body.ClientID,
		ProfessionalID: body.ProfessionalID,
		HourRate:       body.HourRate,
		TotalAmount:    body.TotalAmount,
		Details:        []byte(body.Details),
		WeekDays:       weekDays,
		Shift:          shared.Shift(body.Shift),
		StartTime:      startTime,
		HoursPerDay:    body.HoursPerDay,
		TotalHours:     body.TotalHours,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	respond(w, http.StatusCreated, toResponse(contract))
}

func (h *ContractHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	claims := authmiddleware.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !authmiddleware.IsAdmin(claims) {
		c, err := h.svc.GetByID(r.Context(), id)
		if err != nil {
			http.Error(w, "contract not found", http.StatusNotFound)
			return
		}
		if !h.validator.ValidateSameUserID(r.Context(), claims.Email, c.ClientID) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
	}

	var body updateContractRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	input := ports.UpdateContractInput{ID: id}
	input.Status = body.Status
	input.HourRate = body.HourRate
	input.TotalAmount = body.TotalAmount
	input.HoursPerDay = body.HoursPerDay
	input.TotalHours = body.TotalHours

	if len(body.Details) > 0 {
		input.Details = []byte(body.Details)
	}
	if len(body.WeekDays) > 0 {
		weekDays := shared.SliceToDayOfWeek(body.WeekDays)
		input.WeekDays = weekDays
	}
	if body.Shift != nil {
		s := shared.Shift(*body.Shift)
		input.Shift = &s
	}
	if body.StartTime != nil {
		t, err := time.Parse("15:04", *body.StartTime)
		if err != nil {
			http.Error(w, "start_time must be in HH:MM format", http.StatusBadRequest)
			return
		}
		input.StartTime = &t
	}

	contract, err := h.svc.Update(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	respond(w, http.StatusOK, toResponse(contract))
}

func (h *ContractHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	claims := authmiddleware.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	c, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "contract not found", http.StatusNotFound)
		return
	}

	if !authmiddleware.IsAdmin(claims) && !h.validator.ValidateSameUserID(r.Context(), claims.Email, c.ClientID) {
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
