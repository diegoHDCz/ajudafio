package authpostgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	authdomain "github.com/diegoHDCz/ajudafio/internal/auth/domain"
	"github.com/diegoHDCz/ajudafio/internal/auth/ports"
	userdomain "github.com/diegoHDCz/ajudafio/internal/user/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db      *pgxpool.Pool
	queries *Queries
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db:      db,
		queries: New(db),
	}
}

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	b := u.Bytes
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func stringToUUID(s string) (pgtype.UUID, error) {
	var u pgtype.UUID
	err := u.Scan(s)
	return u, err
}

func tsToTime(t pgtype.Timestamp) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

func timeToTs(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: t, Valid: true}
}

func inferRole(roles []string) userdomain.Role {
	for _, r := range roles {
		switch r {
		case "PROFESSIONAL":
			return userdomain.RoleProfessional
		case "ADMIN":
			return userdomain.RoleAdmin
		}
	}
	return userdomain.RoleClient
}

func toUserDomain(u User) *userdomain.User {
	return &userdomain.User{
		ID:        userdomain.UserID(uuidToString(u.ID)),
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      userdomain.Role(u.Role),
		CreatedAt: tsToTime(u.CreatedAt),
		UpdatedAt: tsToTime(u.UpdatedAt),
	}
}

func toPortAccount(a Account) *ports.Account {
	return &ports.Account{
		ID:                    uuidToString(a.ID),
		AccountID:             a.AccountID,
		ProviderID:            a.ProviderID,
		UserID:                uuidToString(a.UserID),
		AccessToken:           a.AccessToken,
		RefreshToken:          a.RefreshToken,
		IDToken:               a.IDToken,
		AccessTokenExpiresAt:  tsToTime(a.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: tsToTime(a.RefreshTokenExpiresAt),
		Scope:                 a.Scope,
		CreatedAt:             tsToTime(a.CreatedAt),
		UpdatedAt:             tsToTime(a.UpdatedAt),
	}
}

func toPortSession(s Session) *ports.Session {
	return &ports.Session{
		ID:        uuidToString(s.ID),
		ExpiresAt: tsToTime(s.ExpiresAt),
		Token:     s.Token,
		IPAddress: s.IpAddress,
		UserAgent: s.UserAgent,
		UserID:    uuidToString(s.UserID),
		CreatedAt: tsToTime(s.CreatedAt),
		UpdatedAt: tsToTime(s.UpdatedAt),
	}
}

func toPortVerification(v Verification) *ports.Verification {
	return &ports.Verification{
		ID:         uuidToString(v.ID),
		Identifier: v.Identifier,
		Value:      v.Value,
		ExpiresAt:  tsToTime(v.ExpiresAt),
		CreatedAt:  tsToTime(v.CreatedAt),
		UpdatedAt:  tsToTime(v.UpdatedAt),
	}
}

func (r *repository) FindOrCreateUser(ctx context.Context, claims authdomain.Claims) (*userdomain.User, error) {
	u, err := r.queries.GetUserByKeycloakID(ctx, claims.UserID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if err == nil {
		return toUserDomain(u), nil
	}

	// No account link yet — upsert user by email
	u, err = r.queries.UpsertUser(ctx, UpsertUserParams{
		Name:  claims.Name,
		Email: claims.Email,
		Role:  string(inferRole(claims.Roles)),
	})
	if err != nil {
		return nil, err
	}
	return toUserDomain(u), nil
}

func (r *repository) GetAccountByProvider(ctx context.Context, providerID, accountID string) (*ports.Account, error) {
	a, err := r.queries.GetAccountByProvider(ctx, GetAccountByProviderParams{
		ProviderID: providerID,
		AccountID:  accountID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toPortAccount(a), nil
}

func (r *repository) CreateAccount(ctx context.Context, params ports.CreateAccountParams) (*ports.Account, error) {
	userID, err := stringToUUID(params.UserID)
	if err != nil {
		return nil, err
	}
	a, err := r.queries.CreateAccount(ctx, CreateAccountParams{
		AccountID:             params.AccountID,
		ProviderID:            params.ProviderID,
		UserID:                userID,
		AccessToken:           params.AccessToken,
		RefreshToken:          params.RefreshToken,
		IDToken:               params.IDToken,
		AccessTokenExpiresAt:  timeToTs(params.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: timeToTs(params.RefreshTokenExpiresAt),
		Scope:                 params.Scope,
	})
	if err != nil {
		return nil, err
	}
	return toPortAccount(a), nil
}

func (r *repository) GetAccountsByUserID(ctx context.Context, userID string) ([]ports.Account, error) {
	pgID, err := stringToUUID(userID)
	if err != nil {
		return nil, err
	}
	accounts, err := r.queries.GetAccountsByUserID(ctx, pgID)
	if err != nil {
		return nil, err
	}
	result := make([]ports.Account, len(accounts))
	for i, a := range accounts {
		result[i] = *toPortAccount(a)
	}
	return result, nil
}

func (r *repository) CreateSession(ctx context.Context, params ports.CreateSessionParams) (*ports.Session, error) {
	userID, err := stringToUUID(params.UserID)
	if err != nil {
		return nil, err
	}
	s, err := r.queries.CreateSession(ctx, CreateSessionParams{
		ExpiresAt: timeToTs(params.ExpiresAt),
		Token:     params.Token,
		IpAddress: params.IPAddress,
		UserAgent: params.UserAgent,
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}
	return toPortSession(s), nil
}

func (r *repository) GetSessionByToken(ctx context.Context, token string) (*ports.Session, error) {
	s, err := r.queries.GetSessionByToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toPortSession(s), nil
}

func (r *repository) DeleteSession(ctx context.Context, token string) error {
	return r.queries.DeleteSession(ctx, token)
}

func (r *repository) DeleteExpiredSessions(ctx context.Context) error {
	return r.queries.DeleteExpiredSessions(ctx)
}

func (r *repository) CreateVerification(ctx context.Context, params ports.CreateVerificationParams) (*ports.Verification, error) {
	v, err := r.queries.CreateVerification(ctx, CreateVerificationParams{
		Identifier: params.Identifier,
		Value:      params.Value,
		ExpiresAt:  timeToTs(params.ExpiresAt),
	})
	if err != nil {
		return nil, err
	}
	return toPortVerification(v), nil
}

func (r *repository) GetVerification(ctx context.Context, identifier, value string) (*ports.Verification, error) {
	v, err := r.queries.GetVerification(ctx, GetVerificationParams{
		Identifier: identifier,
		Value:      value,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toPortVerification(v), nil
}

func (r *repository) DeleteVerification(ctx context.Context, identifier, value string) error {
	return r.queries.DeleteVerification(ctx, DeleteVerificationParams{
		Identifier: identifier,
		Value:      value,
	})
}
