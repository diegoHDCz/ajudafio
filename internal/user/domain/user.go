package domain

import "time"

type Role string

const (
	RoleFamilyClient       Role = "FAMILY_CLIENT"
	RoleHealthCareprovider Role = "HEALTH_CAREPROVIDER"
	RoleFinancialSponsor   Role = "FINANCIAL_SPONSOR"
	RolePlatformAdmin      Role = "PLATFORM_ADMIN"
)

type OnboardingStatus string

const (
	OnboardingNotStarted           OnboardingStatus = "NOT_STARTED"
	OnboardingInProgress           OnboardingStatus = "IN_PROGRESS"
	OnboardingDocumentationPending OnboardingStatus = "DOCUMENTATION_PENDING"
	OnboardingUnderReview          OnboardingStatus = "UNDER_REVIEW"
	OnboardingCompleted            OnboardingStatus = "COMPLETED"
)

type User struct {
	ID               string
	AuthUserID       *string // Supabase auth.users.id — nil until the user is provisioned via Supabase Auth.
	Name             string  // Alterado para string (NOT NULL no banco)
	Email            string  // NOT NULL UNIQUE no banco
	Phone            *string // Ponteiro pois na tabela não tem NOT NULL
	Role             Role
	OnboardingStatus OnboardingStatus
	AvatarURL        *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
