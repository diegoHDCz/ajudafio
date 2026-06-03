package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/diegoHDCz/ajudafio/internal/appointment/ports"
	"github.com/go-chi/chi/v5"
)

type AppointmentHandler struct {
	svc ports.AppointmentService
}

func NewHandler(svc ports.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{svc: svc}
}

func NewRouter(h *AppointmentHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Get("/contract/{contractID}", h.GetByContractID)
	r.Get("/client/{clientID}", h.GetByClientID)
	r.Get("/professional/{professionalID}", h.GetByProfessionalID)
	r.Patch("/{id}/status", h.UpdateStatus)
	r.Delete("/{id}", h.Delete)
	return r
}

// @Summary      Criar agendamento
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        body  body      createAppointmentRequest  true  "Dados do agendamento"
// @Success      201   {object}  appointmentResponse
// @Failure      400   {string}  string
// @Failure      422   {string}  string
// @Security     BearerAuth
// @Router       /appointments [post]
func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body createAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.ContractID == "" || body.ClientID == "" || body.ProfessionalID == "" {
		http.Error(w, "contract_id, client_id and professional_id are required", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", body.Date)
	if err != nil {
		http.Error(w, "date must be in YYYY-MM-DD format", http.StatusBadRequest)
		return
	}
	startTime, err := time.Parse("15:04", body.StartTime)
	if err != nil {
		http.Error(w, "start_time must be in HH:MM format", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse("15:04", body.EndTime)
	if err != nil {
		http.Error(w, "end_time must be in HH:MM format", http.StatusBadRequest)
		return
	}

	a, err := h.svc.Create(r.Context(), ports.CreateAppointmentInput{
		ContractID:     body.ContractID,
		ClientID:       body.ClientID,
		ProfessionalID: body.ProfessionalID,
		Date:           date,
		StartTime:      startTime,
		EndTime:        endTime,
		ZipCode:        body.ZipCode,
		AddressLine:    body.AddressLine,
		Number:         body.Number,
		Complement:     body.Complement,
		District:       body.District,
		City:           body.City,
		State:          body.State,
		Reference:      body.Reference,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	respond(w, http.StatusCreated, toResponse(a))
}

// @Summary      Buscar agendamento por ID
// @Tags         appointments
// @Produce      json
// @Param        id  path      string  true  "Appointment ID"
// @Success      200  {object}  appointmentResponse
// @Failure      404  {string}  string
// @Security     BearerAuth
// @Router       /appointments/{id} [get]
func (h *AppointmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	a, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "appointment not found", http.StatusNotFound)
		return
	}
	respond(w, http.StatusOK, toResponse(a))
}

// @Summary      Buscar agendamentos por contrato
// @Tags         appointments
// @Produce      json
// @Param        contractID  path      string  true  "Contract ID"
// @Success      200         {array}   appointmentResponse
// @Failure      500         {string}  string
// @Security     BearerAuth
// @Router       /appointments/contract/{contractID} [get]
func (h *AppointmentHandler) GetByContractID(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "contractID")
	list, err := h.svc.GetByContractID(r.Context(), contractID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]appointmentResponse, len(list))
	for i, a := range list {
		resp[i] = toResponse(a)
	}
	respond(w, http.StatusOK, resp)
}

// @Summary      Buscar agendamentos por cliente
// @Tags         appointments
// @Produce      json
// @Param        clientID  path      string  true  "Client ID"
// @Success      200       {array}   appointmentResponse
// @Failure      500       {string}  string
// @Security     BearerAuth
// @Router       /appointments/client/{clientID} [get]
func (h *AppointmentHandler) GetByClientID(w http.ResponseWriter, r *http.Request) {
	clientID := chi.URLParam(r, "clientID")
	list, err := h.svc.GetByClientID(r.Context(), clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]appointmentResponse, len(list))
	for i, a := range list {
		resp[i] = toResponse(a)
	}
	respond(w, http.StatusOK, resp)
}

// @Summary      Buscar agendamentos por profissional
// @Tags         appointments
// @Produce      json
// @Param        professionalID  path      string  true  "Professional ID"
// @Success      200             {array}   appointmentResponse
// @Failure      500             {string}  string
// @Security     BearerAuth
// @Router       /appointments/professional/{professionalID} [get]
func (h *AppointmentHandler) GetByProfessionalID(w http.ResponseWriter, r *http.Request) {
	professionalID := chi.URLParam(r, "professionalID")
	list, err := h.svc.GetByProfessionalID(r.Context(), professionalID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]appointmentResponse, len(list))
	for i, a := range list {
		resp[i] = toResponse(a)
	}
	respond(w, http.StatusOK, resp)
}

// @Summary      Atualizar status do agendamento
// @Tags         appointments
// @Accept       json
// @Produce      json
// @Param        id    path      string               true  "Appointment ID"
// @Param        body  body      updateStatusRequest  true  "Novo status"
// @Success      200   {object}  appointmentResponse
// @Failure      400   {string}  string
// @Failure      404   {string}  string
// @Failure      422   {string}  string
// @Security     BearerAuth
// @Router       /appointments/{id}/status [patch]
func (h *AppointmentHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var body updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Status == "" {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}

	a, err := h.svc.UpdateStatus(r.Context(), id, body.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	respond(w, http.StatusOK, toResponse(a))
}

// @Summary      Remover agendamento
// @Tags         appointments
// @Param        id  path  string  true  "Appointment ID"
// @Success      204
// @Failure      500  {string}  string
// @Security     BearerAuth
// @Router       /appointments/{id} [delete]
func (h *AppointmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
