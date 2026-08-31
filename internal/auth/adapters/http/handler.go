package http

import (
	"encoding/json"
	"net/http"

	auditdomain "github.com/diegoHDCz/ajudafio/internal/audit/domain"
	auditports "github.com/diegoHDCz/ajudafio/internal/audit/ports"
	"github.com/diegoHDCz/ajudafio/internal/auth/middleware"
	"github.com/diegoHDCz/ajudafio/internal/auth/rbac"
	profileports "github.com/diegoHDCz/ajudafio/internal/profile/ports"
	userdomain "github.com/diegoHDCz/ajudafio/internal/user/domain"
	userports "github.com/diegoHDCz/ajudafio/internal/user/ports"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	userSvc    userports.UserService
	profileSvc profileports.ProfileService
	auditSvc   auditports.AuditService
}

func NewHandler(userSvc userports.UserService, profileSvc profileports.ProfileService, auditSvc auditports.AuditService) *Handler {
	return &Handler{userSvc: userSvc, profileSvc: profileSvc, auditSvc: auditSvc}
}

// NewRouter mounts the endpoints still owned by this module now that
// Supabase Auth handles login/registration directly with the client (see
// ADR-005). Everything under /me and /onboarding is wired at the top level in
// cmd/main.go, matching the flat paths in docs/feat/refactor-auth.md §32.
func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Post("/profile", h.CompleteProfile)
	return r
}

// @Summary      Dados do usuário autenticado (cria o usuário de aplicação no primeiro acesso)
// @Tags         auth
// @Produce      json
// @Success      200  {object}  meResponse
// @Failure      401  {string}  string
// @Security     BearerAuth
// @Router       /me [get]
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if claims.UserID != "" {
		user, err := h.userSvc.GetByID(r.Context(), claims.UserID)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		respond(w, http.StatusOK, toMeResponse(user))
		return
	}

	// First access (feat doc §24): an authenticated Supabase identity with no
	// application user yet — provision one instead of 404ing.
	user, err := h.userSvc.EnsureProvisioned(r.Context(), userports.ProvisionUserInput{
		AuthUserID: claims.AuthUserID,
		Email:      claims.Email,
		Name:       claims.Email,
	})
	if err != nil {
		http.Error(w, "failed to provision user", http.StatusInternalServerError)
		return
	}
	h.auditSvc.Log(r.Context(), &user.ID, auditdomain.ActionUserRegistered, map[string]any{"auth_user_id": claims.AuthUserID})

	respond(w, http.StatusOK, toMeResponse(user))
}

// @Summary      Papéis do usuário autenticado
// @Tags         auth
// @Produce      json
// @Success      200  {object}  rolesResponse
// @Failure      401  {string}  string
// @Failure      403  {string}  string
// @Security     BearerAuth
// @Router       /me/roles [get]
func (h *Handler) MeRoles(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil || claims.UserID == "" {
		http.Error(w, "user not initialized", http.StatusForbidden)
		return
	}

	roles, err := h.userSvc.ListRoles(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, "failed to load roles", http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusOK, rolesResponse{Roles: rolesToStrings(roles)})
}

// @Summary      Permissões do usuário autenticado
// @Tags         auth
// @Produce      json
// @Success      200  {object}  permissionsResponse
// @Failure      401  {string}  string
// @Failure      403  {string}  string
// @Security     BearerAuth
// @Router       /me/permissions [get]
func (h *Handler) MePermissions(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil || claims.UserID == "" {
		http.Error(w, "user not initialized", http.StatusForbidden)
		return
	}

	roles, err := h.userSvc.ListRoles(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, "failed to load roles", http.StatusInternalServerError)
		return
	}

	perms := rbac.PermissionsForRoles(rolesToStrings(roles))
	permStrings := make([]string, len(perms))
	for i, p := range perms {
		permStrings[i] = string(p)
	}

	respond(w, http.StatusOK, permissionsResponse{Permissions: permStrings})
}

// @Summary      Completar perfil (FAMILY_CLIENT ou FINANCIAL_SPONSOR)
// @Tags         auth
// @Accept       json
// @Param        body  body  completeProfileRequest  true  "Papel escolhido"
// @Success      204
// @Failure      400  {string}  string
// @Failure      401  {string}  string
// @Failure      403  {string}  string
// @Security     BearerAuth
// @Router       /auth/profile [post]
func (h *Handler) CompleteProfile(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil || claims.UserID == "" {
		http.Error(w, "user not initialized", http.StatusForbidden)
		return
	}

	var body completeProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	switch userdomain.Role(body.Role) {
	case userdomain.RoleFamilyClient:
		if _, err := h.profileSvc.CreateFamilyProfile(r.Context(), claims.UserID); err != nil {
			http.Error(w, "failed to create profile", http.StatusInternalServerError)
			return
		}
	case userdomain.RoleFinancialSponsor:
		if _, err := h.profileSvc.CreateFinancialProfile(r.Context(), claims.UserID); err != nil {
			http.Error(w, "failed to create profile", http.StatusInternalServerError)
			return
		}
	case userdomain.RoleHealthCareprovider:
		http.Error(w, "use POST /professionals to create a health care provider profile", http.StatusBadRequest)
		return
	default:
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	if err := h.userSvc.UpdateUserRole(r.Context(), claims.UserID, userdomain.Role(body.Role)); err != nil {
		http.Error(w, "failed to update role", http.StatusInternalServerError)
		return
	}
	h.auditSvc.Log(r.Context(), &claims.UserID, auditdomain.ActionRoleAssigned, map[string]any{"role": body.Role})

	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Atualizar status de onboarding
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      onboardingStatusRequest  true  "Novo status"
// @Success      200   {object}  onboardingStatusResponse
// @Failure      400   {string}  string
// @Failure      401   {string}  string
// @Failure      403   {string}  string
// @Security     BearerAuth
// @Router       /onboarding [post]
func (h *Handler) UpdateOnboarding(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil || claims.UserID == "" {
		http.Error(w, "user not initialized", http.StatusForbidden)
		return
	}

	var body onboardingStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	user, err := h.userSvc.UpdateOnboardingStatus(r.Context(), claims.UserID, userdomain.OnboardingStatus(body.Status))
	if err != nil {
		http.Error(w, "failed to update onboarding status", http.StatusInternalServerError)
		return
	}

	respond(w, http.StatusOK, onboardingStatusResponse{Status: string(user.OnboardingStatus)})
}

// @Summary      Status de onboarding do usuário autenticado
// @Tags         auth
// @Produce      json
// @Success      200  {object}  onboardingStatusResponse
// @Failure      401  {string}  string
// @Failure      403  {string}  string
// @Security     BearerAuth
// @Router       /onboarding/status [get]
func (h *Handler) OnboardingStatus(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil || claims.UserID == "" {
		http.Error(w, "user not initialized", http.StatusForbidden)
		return
	}

	user, err := h.userSvc.GetByID(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	respond(w, http.StatusOK, onboardingStatusResponse{Status: string(user.OnboardingStatus)})
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func rolesToStrings(roles []userdomain.Role) []string {
	out := make([]string, len(roles))
	for i, r := range roles {
		out[i] = string(r)
	}
	return out
}
