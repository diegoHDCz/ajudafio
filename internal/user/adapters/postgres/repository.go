package userpostgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/diegoHDCz/ajudafio/internal/user/domain"
	"github.com/diegoHDCz/ajudafio/internal/user/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewRepository(db *pgxpool.Pool) ports.UserRepository {
	return &repository{
		db:      db,
		queries: New(db),
	}
}

func (r *repository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("userpostgres.GetByID: invalid id: %w", err)
	}
	row, err := r.queries.GetUserByID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("userpostgres.GetByID: %w", err)
	}
	return mapUser(row.ID, row.AuthUserID, row.Name, row.Email, row.Phone, row.Role, row.OnboardingStatus, row.AvatarUrl, row.CreatedAt, row.UpdatedAt), nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("userpostgres.GetByEmail: %w", err)
	}
	return mapUser(row.ID, row.AuthUserID, row.Name, row.Email, row.Phone, row.Role, row.OnboardingStatus, row.AvatarUrl, row.CreatedAt, row.UpdatedAt), nil
}

func (r *repository) GetByAuthUserID(ctx context.Context, authUserID string) (*domain.User, error) {
	aid, err := parseUUID(authUserID)
	if err != nil {
		return nil, fmt.Errorf("userpostgres.GetByAuthUserID: invalid auth user id: %w", err)
	}
	row, err := r.queries.GetUserByAuthUserID(ctx, aid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("userpostgres.GetByAuthUserID: %w", err)
	}
	return mapUser(row.ID, row.AuthUserID, row.Name, row.Email, row.Phone, row.Role, row.OnboardingStatus, row.AvatarUrl, row.CreatedAt, row.UpdatedAt), nil
}

func (r *repository) FindOrCreateByAuthUserID(ctx context.Context, authUserID, email, name string, defaultRole domain.Role) (*domain.User, error) {
	aid, err := parseUUID(authUserID)
	if err != nil {
		return nil, fmt.Errorf("userpostgres.FindOrCreateByAuthUserID: invalid auth user id: %w", err)
	}
	row, err := r.queries.FindOrCreateUserByAuthUserID(ctx, FindOrCreateUserByAuthUserIDParams{
		ID:         pgtype.UUID{Bytes: uuid.New(), Valid: true},
		AuthUserID: aid,
		Name:       name,
		Email:      email,
		Role:       string(defaultRole),
	})
	if err != nil {
		return nil, fmt.Errorf("userpostgres.FindOrCreateByAuthUserID: %w", err)
	}
	return mapUser(row.ID, row.AuthUserID, row.Name, row.Email, row.Phone, row.Role, row.OnboardingStatus, row.AvatarUrl, row.CreatedAt, row.UpdatedAt), nil
}

func (r *repository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	uid, err := parseUUID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("userpostgres.Create: invalid id: %w", err)
	}
	var authUID pgtype.UUID
	if user.AuthUserID != nil {
		authUID, err = parseUUID(*user.AuthUserID)
		if err != nil {
			return nil, fmt.Errorf("userpostgres.Create: invalid auth user id: %w", err)
		}
	}
	row, err := r.queries.CreateUser(ctx, CreateUserParams{
		ID:         uid,
		AuthUserID: authUID,
		Name:       user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		Role:       string(user.Role),
	})
	if err != nil {
		return nil, fmt.Errorf("userpostgres.Create: %w", err)
	}
	return mapUser(row.ID, row.AuthUserID, row.Name, row.Email, row.Phone, row.Role, row.OnboardingStatus, row.AvatarUrl, row.CreatedAt, row.UpdatedAt), nil
}

func (r *repository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	uid, err := parseUUID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("userpostgres.Update: invalid id: %w", err)
	}
	row, err := r.queries.UpdateUser(ctx, UpdateUserParams{
		ID:    uid,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Role:  string(user.Role),
	})
	if err != nil {
		return nil, fmt.Errorf("userpostgres.Update: %w", err)
	}
	return mapUser(row.ID, row.AuthUserID, row.Name, row.Email, row.Phone, row.Role, row.OnboardingStatus, row.AvatarUrl, row.CreatedAt, row.UpdatedAt), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	uid, err := parseUUID(id)
	if err != nil {
		return fmt.Errorf("userpostgres.Delete: invalid id: %w", err)
	}
	if err := r.queries.DeleteUser(ctx, uid); err != nil {
		return fmt.Errorf("userpostgres.Delete: %w", err)
	}
	return nil
}

func (r *repository) UpdateUserRole(ctx context.Context, id string, role domain.Role) error {
	uid, err := parseUUID(id)
	if err != nil {
		return fmt.Errorf("userpostgres.UpdateUserRole: invalid id: %w", err)
	}
	if err := r.queries.UpdateUserRole(ctx, UpdateUserRoleParams{
		ID:   uid,
		Role: string(role),
	}); err != nil {
		return fmt.Errorf("userpostgres.UpdateUserRole: %w", err)
	}
	return nil
}

func (r *repository) UpdateAvatar(ctx context.Context, id string, avatarURL *string) (*domain.User, error) {
	uid, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("userpostgres.UpdateAvatar: invalid id: %w", err)
	}
	row, err := r.queries.UpdateUserAvatar(ctx, UpdateUserAvatarParams{
		ID:        uid,
		AvatarUrl: avatarURL,
	})
	if err != nil {
		return nil, fmt.Errorf("userpostgres.UpdateAvatar: %w", err)
	}
	return mapUser(row.ID, row.AuthUserID, row.Name, row.Email, row.Phone, row.Role, row.OnboardingStatus, row.AvatarUrl, row.CreatedAt, row.UpdatedAt), nil
}

func (r *repository) UpdateOnboardingStatus(ctx context.Context, id string, status domain.OnboardingStatus) error {
	uid, err := parseUUID(id)
	if err != nil {
		return fmt.Errorf("userpostgres.UpdateOnboardingStatus: invalid id: %w", err)
	}
	if err := r.queries.UpdateOnboardingStatus(ctx, UpdateOnboardingStatusParams{
		ID:               uid,
		OnboardingStatus: string(status),
	}); err != nil {
		return fmt.Errorf("userpostgres.UpdateOnboardingStatus: %w", err)
	}
	return nil
}

func (r *repository) ListRoles(ctx context.Context, userID string) ([]domain.Role, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("userpostgres.ListRoles: invalid id: %w", err)
	}
	rows, err := r.queries.ListRolesByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("userpostgres.ListRoles: %w", err)
	}
	roles := make([]domain.Role, 0, len(rows))
	for _, row := range rows {
		roles = append(roles, domain.Role(row.Role))
	}
	return roles, nil
}

func (r *repository) AddRole(ctx context.Context, userID string, role domain.Role) error {
	uid, err := parseUUID(userID)
	if err != nil {
		return fmt.Errorf("userpostgres.AddRole: invalid id: %w", err)
	}
	// ON CONFLICT DO NOTHING means no row (pgx.ErrNoRows) when the role already exists — not an error.
	if _, err := r.queries.AddUserRole(ctx, AddUserRoleParams{UserID: uid, Role: string(role)}); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("userpostgres.AddRole: %w", err)
	}
	return nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func mapUser(
	id, authUserID pgtype.UUID,
	name, email string,
	phone *string,
	role string,
	onboardingStatus string,
	avatarURL *string,
	createdAt, updatedAt pgtype.Timestamp,
) *domain.User {
	var authUserIDPtr *string
	if authUserID.Valid {
		s := uuid.UUID(authUserID.Bytes).String()
		authUserIDPtr = &s
	}
	return &domain.User{
		ID:               uuid.UUID(id.Bytes).String(),
		AuthUserID:       authUserIDPtr,
		Name:             name,
		Email:            email,
		Phone:            phone,
		Role:             domain.Role(role),
		OnboardingStatus: domain.OnboardingStatus(onboardingStatus),
		AvatarURL:        avatarURL,
		CreatedAt:        createdAt.Time,
		UpdatedAt:        updatedAt.Time,
	}
}

func parseUUID(s string) (pgtype.UUID, error) {
	uid, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: uid, Valid: true}, nil
}
