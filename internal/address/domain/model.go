package domain

import "time"

// Address representa a estrutura da tabela 'addresses'
type Address struct {
	ID          string
	UserID      string
	ContractID  *string // Ponteiro para aceitar NULL
	ZipCode     string
	AddressLine string
	Number      string
	Complement  *string // Ponteiro para aceitar NULL
	District    string
	City        string
	State       string
	Reference   *string // Ponteiro para aceitar NULL
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
