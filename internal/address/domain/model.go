package domain

import "time"

// Address representa a estrutura da tabela 'addresses'
type Address struct {
	ID          string
	UserID      string
	ZipCode     string
	AddressLine string
	Number      string
	Complement  *string
	District    string
	City        string
	State       string
	Reference   *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
