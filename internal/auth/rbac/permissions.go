// Package rbac maps app roles to permissions (feat doc §21). It operates on
// plain role strings rather than internal/user/domain.Role to avoid coupling
// the auth module to the user module's domain package.
package rbac

type Permission string

const (
	PermViewOwnProfile       Permission = "VIEW_OWN_PROFILE"
	PermUpdateOwnProfile     Permission = "UPDATE_OWN_PROFILE"
	PermCreateServiceRequest Permission = "CREATE_SERVICE_REQUEST"
	PermViewServiceRequest   Permission = "VIEW_SERVICE_REQUEST"
	PermChatProvider         Permission = "CHAT_PROVIDER"

	PermUpdateAvailability   Permission = "UPDATE_AVAILABILITY"
	PermAcceptServiceRequest Permission = "ACCEPT_SERVICE_REQUEST"
	PermChatFamily           Permission = "CHAT_FAMILY"

	PermViewBilling      Permission = "VIEW_BILLING"
	PermViewInvoices     Permission = "VIEW_INVOICES"
	PermMakePayment      Permission = "MAKE_PAYMENT"
	PermViewTransactions Permission = "VIEW_TRANSACTIONS"

	PermManageUsers     Permission = "MANAGE_USERS"
	PermManageProviders Permission = "MANAGE_PROVIDERS"
	PermManageRoles     Permission = "MANAGE_ROLES"
	PermViewAuditLog    Permission = "VIEW_AUDIT_LOG"
	PermManagePlatform  Permission = "MANAGE_PLATFORM"
)

const (
	RoleFamilyClient       = "FAMILY_CLIENT"
	RoleHealthCareprovider = "HEALTH_CAREPROVIDER"
	RoleFinancialSponsor   = "FINANCIAL_SPONSOR"
	RolePlatformAdmin      = "PLATFORM_ADMIN"
)

var rolePermissions = map[string][]Permission{
	RoleFamilyClient: {
		PermViewOwnProfile, PermUpdateOwnProfile, PermCreateServiceRequest, PermViewServiceRequest, PermChatProvider,
	},
	RoleHealthCareprovider: {
		PermViewOwnProfile, PermUpdateOwnProfile, PermUpdateAvailability, PermViewServiceRequest, PermAcceptServiceRequest, PermChatFamily,
	},
	RoleFinancialSponsor: {
		PermViewBilling, PermViewInvoices, PermMakePayment, PermViewTransactions,
	},
	RolePlatformAdmin: {
		PermManageUsers, PermManageProviders, PermManageRoles, PermViewAuditLog, PermManagePlatform,
	},
}

// PermissionsForRoles returns the de-duplicated union of permissions granted
// to any of the given roles, in a stable order.
func PermissionsForRoles(roles []string) []Permission {
	seen := make(map[Permission]bool)
	var result []Permission
	for _, role := range roles {
		for _, perm := range rolePermissions[role] {
			if !seen[perm] {
				seen[perm] = true
				result = append(result, perm)
			}
		}
	}
	return result
}

// HasPermission reports whether any of the given roles grants perm. Roles
// must come from a source resolved server-side (users/user_roles) — never
// from a client-supplied value (feat doc §22).
func HasPermission(roles []string, perm Permission) bool {
	for _, role := range roles {
		for _, p := range rolePermissions[role] {
			if p == perm {
				return true
			}
		}
	}
	return false
}
