package domain

import "time"

type Log struct {
	ID        string
	UserID    *string
	Action    string
	Metadata  map[string]any
	CreatedAt time.Time
}

// Action names per docs/feat/refactor-auth.md §31.
const (
	ActionUserRegistered         = "USER_REGISTERED"
	ActionUserEmailVerified      = "USER_EMAIL_VERIFIED"
	ActionUserLogin              = "USER_LOGIN"
	ActionUserLogout             = "USER_LOGOUT"
	ActionRoleAssigned           = "ROLE_ASSIGNED"
	ActionRoleChanged            = "ROLE_CHANGED"
	ActionUserBlocked            = "USER_BLOCKED"
	ActionUserUnblocked          = "USER_UNBLOCKED"
	ActionProviderApproved       = "PROVIDER_APPROVED"
	ActionProviderRejected       = "PROVIDER_REJECTED"
	ActionPasswordResetRequested = "PASSWORD_RESET_REQUESTED"
)
