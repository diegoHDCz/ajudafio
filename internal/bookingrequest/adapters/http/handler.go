package http

import (
	"encoding/json"
	"net/http"

	authmiddleware "github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	"github.com/diegoHDCz/ajudafio/internal/bookingrequest/domain"
	"github.com/diegoHDCz/ajudafio/internal/bookingrequest/ports"
	professionalports "github.com/diegoHDCz/ajudafio/internal/professional/ports"
	"github.com/diegoHDCz/ajudafio/internal/shared"
	"github.com/go-chi/chi/v5"
)

type BookingRequestHandler struct {
	svc             ports.BookingRequestService
	validator       *shared.Validator
	professionalSvc professionalports.ProfessionalService
}

func NewHandler(svc ports.BookingRequestService, validator *shared.Validator, professionalSvc professionalports.ProfessionalService) *BookingRequestHandler {
	return &BookingRequestHandler{svc: svc, validator: validator, professionalSvc: professionalSvc}
}

func NewRouter(h *BookingRequestHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Get("/client/{clientID}", h.GetByClientID)
	r.Get("/professional/{professionalID}", h.GetByProfessionalID)
	r.Patch("/{id}/status", h.UpdateStatus)
	return r
}

// @Summary      Criar solicitação de agendamento
// @Tags         booking-requests
// @Accept       json
// @Produce      json
// @Param        body  body      createBookingRequestRequest  true  "Dados da solicitação"
// @Success      201   {object}  bookingRequestResponse
// @Failure      400   {string}  string
// @Failure      422   {string}  string
// @Security     BearerAuth
// @Router       /booking-requests [post]
func (h *BookingRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body createBookingRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.ClientID == "" || body.ProfessionalID == "" || body.AddressID == "" {
		http.Error(w, "client_id, professional_id and address_id are required", http.StatusBadRequest)
		return
	}
	if body.ProposedValue <= 0 {
		http.Error(w, "proposed_value must be greater than zero", http.StatusBadRequest)
		return
	}

	schedule := domain.ScheduleDetails{
		StartDate: body.ScheduleDetails.StartDate,
		EndDate:   body.ScheduleDetails.EndDate,
	}
	for _, e := range body.ScheduleDetails.Recurrence {
		schedule.Recurrence = append(schedule.Recurrence, domain.ScheduleEntry{
			DayOfWeek: e.DayOfWeek,
			StartTime: e.StartTime,
			EndTime:   e.EndTime,
		})
	}

	br, err := h.svc.Create(r.Context(), ports.CreateBookingRequestInput{
		ClientID:        body.ClientID,
		ProfessionalID:  body.ProfessionalID,
		AddressID:       body.AddressID,
		ProposedValue:   body.ProposedValue,
		ScheduleDetails: schedule,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	respond(w, http.StatusCreated, toResponse(br))
}

// @Summary      Buscar solicitação por ID
// @Tags         booking-requests
// @Produce      json
// @Param        id  path      string  true  "Booking Request ID"
// @Success      200  {object}  bookingRequestResponse
// @Failure      404  {string}  string
// @Security     BearerAuth
// @Router       /booking-requests/{id} [get]
func (h *BookingRequestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	br, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "booking request not found", http.StatusNotFound)
		return
	}
	respond(w, http.StatusOK, toResponse(br))
}

// @Summary      Buscar solicitações por cliente
// @Tags         booking-requests
// @Produce      json
// @Param        clientID  path      string  true  "Client ID"
// @Success      200       {array}   bookingRequestResponse
// @Failure      500       {string}  string
// @Security     BearerAuth
// @Router       /booking-requests/client/{clientID} [get]
func (h *BookingRequestHandler) GetByClientID(w http.ResponseWriter, r *http.Request) {
	clientID := chi.URLParam(r, "clientID")
	list, err := h.svc.GetByClientID(r.Context(), clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]bookingRequestResponse, len(list))
	for i, br := range list {
		resp[i] = toResponse(br)
	}
	respond(w, http.StatusOK, resp)
}

// @Summary      Buscar solicitações por profissional
// @Tags         booking-requests
// @Produce      json
// @Param        professionalID  path      string  true  "Professional ID"
// @Success      200             {array}   bookingRequestResponse
// @Failure      500             {string}  string
// @Security     BearerAuth
// @Router       /booking-requests/professional/{professionalID} [get]
func (h *BookingRequestHandler) GetByProfessionalID(w http.ResponseWriter, r *http.Request) {
	professionalID := chi.URLParam(r, "professionalID")
	list, err := h.svc.GetByProfessionalID(r.Context(), professionalID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]bookingRequestResponse, len(list))
	for i, br := range list {
		resp[i] = toResponse(br)
	}
	respond(w, http.StatusOK, resp)
}

// @Summary      Atualizar status da solicitação
// @Tags         booking-requests
// @Accept       json
// @Produce      json
// @Param        id    path      string               true  "Booking Request ID"
// @Param        body  body      updateStatusRequest  true  "Novo status"
// @Success      200   {object}  bookingRequestResponse
// @Failure      400   {string}  string
// @Failure      401   {string}  string
// @Failure      403   {string}  string
// @Failure      404   {string}  string
// @Failure      422   {string}  string
// @Security     BearerAuth
// @Router       /booking-requests/{id}/status [patch]
func (h *BookingRequestHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	claims := authmiddleware.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var body updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Status == "" {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}

	if !authmiddleware.IsAdmin(claims) {
		if err := h.checkStatusOwnership(r, id, body.Status); err != nil {
			http.Error(w, err.Error(), err.(httpErr).status)
			return
		}
	}

	br, err := h.svc.UpdateStatus(r.Context(), id, body.Status, body.RejectionReason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	respond(w, http.StatusOK, toResponse(br))
}

// checkStatusOwnership validates that the caller is allowed to apply the given status transition.
// ACCEPTED/REJECTED → only the professional; CANCELLED → only the client.
func (h *BookingRequestHandler) checkStatusOwnership(r *http.Request, id, status string) error {
	claims := authmiddleware.GetClaims(r.Context())

	br, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		return httpErr{"booking request not found", http.StatusNotFound}
	}

	switch domain.BookingRequestStatus(status) {
	case domain.StatusAccepted, domain.StatusRejected:
		p, err := h.professionalSvc.GetByID(r.Context(), br.ProfessionalID)
		if err != nil {
			return httpErr{"professional not found", http.StatusNotFound}
		}
		if !h.validator.ValidateSameUserID(r.Context(), claims.Email, p.UserID) {
			return httpErr{"forbidden", http.StatusForbidden}
		}
	case domain.StatusCancelled:
		if !h.validator.ValidateSameUserID(r.Context(), claims.Email, br.ClientID) {
			return httpErr{"forbidden", http.StatusForbidden}
		}
	}
	return nil
}

type httpErr struct {
	msg    string
	status int
}

func (e httpErr) Error() string { return e.msg }

func toResponse(br *domain.BookingRequest) bookingRequestResponse {
	entries := make([]scheduleEntryResponse, len(br.ScheduleDetails.Recurrence))
	for i, e := range br.ScheduleDetails.Recurrence {
		entries[i] = scheduleEntryResponse{
			DayOfWeek: e.DayOfWeek,
			StartTime: e.StartTime,
			EndTime:   e.EndTime,
		}
	}
	return bookingRequestResponse{
		ID:             br.ID,
		ClientID:       br.ClientID,
		ProfessionalID: br.ProfessionalID,
		AddressID:      br.AddressID,
		ProposedValue:  br.ProposedValue,
		ScheduleDetails: scheduleDetailsResponse{
			Recurrence: entries,
			StartDate:  br.ScheduleDetails.StartDate,
			EndDate:    br.ScheduleDetails.EndDate,
		},
		Status:          string(br.Status),
		RejectionReason: br.RejectionReason,
		CreatedAt:       br.CreatedAt,
		RespondedAt:     br.RespondedAt,
	}
}

func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
