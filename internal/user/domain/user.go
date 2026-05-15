package domain

import "time"

// UserID definido como string para representar o UUID vindo do banco.
type UserID string

type Role string

const (
	RoleClient       Role = "CLIENT"
	RoleProfessional Role = "PROFESSIONAL"
	RoleAdmin        Role = "ADMIN"
)

type User struct {
	ID        UserID
	Name      string  // Alterado para string (NOT NULL no banco)
	Email     string  // NOT NULL UNIQUE no banco
	Phone     *string // Ponteiro pois na tabela não tem NOT NULL
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}
