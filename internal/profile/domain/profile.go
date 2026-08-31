package domain

import "time"

// FamilyProfile and FinancialProfile are intentionally minimal — the feat doc
// (docs/feat/refactor-auth.md §8) doesn't specify concrete fields for these
// roles yet, unlike HEALTH_CAREPROVIDER's profile, which already exists as
// internal/professional and is reused as-is rather than duplicated here.

type FamilyProfile struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FinancialProfile struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
