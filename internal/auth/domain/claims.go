package domain

// AuthenticatedUser is the resolved identity attached to the request context
// once a Supabase-issued JWT has been validated and matched to a local user
// (or, for a first-time caller, matched only by AuthUserID with UserID empty —
// see the /me provisioning flow).
type AuthenticatedUser struct {
	AuthUserID string // Supabase auth.users.id — the JWT `sub` claim.
	UserID     string // Local users.id, empty until the user is provisioned.
	Email      string
	Name       string
	Role       string // App RBAC role, resolved from users/user_roles — never the raw Supabase JWT claim.
}
