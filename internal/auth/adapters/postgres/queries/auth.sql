-- =============================================================================
-- Users
-- =============================================================================

-- name: UpsertUser :one
-- Upserts a local user from Keycloak claims forwarded by KrakenD.
-- On conflict by email, refreshes the name to reflect Keycloak profile changes.
INSERT INTO users (name, email, role, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT (email) DO UPDATE
    SET name       = EXCLUDED.name,
        updated_at = NOW()
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByKeycloakID :one
-- Resolves a local user via the Keycloak subject (sub claim) stored in accounts.
SELECT
    u.id, u.name, u.email, u.phone, u.role, u.created_at, u.updated_at
FROM users u
JOIN accounts a ON a.user_id = u.id
WHERE a.provider_id = 'keycloak' AND a.account_id = $1
LIMIT 1;

-- =============================================================================
-- Accounts (provider links)
-- =============================================================================

-- name: CreateAccount :one
-- Links a Keycloak account to a local user. password is always NULL for OIDC providers.
INSERT INTO accounts (
    account_id, provider_id, user_id,
    access_token, refresh_token, id_token,
    access_token_expires_at, refresh_token_expires_at,
    scope, password, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, NULL, NOW(), NOW()
)
RETURNING *;

-- name: GetAccountByProvider :one
SELECT * FROM accounts
WHERE provider_id = $1 AND account_id = $2 LIMIT 1;

-- name: GetAccountsByUserID :many
SELECT * FROM accounts
WHERE user_id = $1;

-- =============================================================================
-- Sessions
-- =============================================================================

-- name: CreateSession :one
INSERT INTO sessions (
    expires_at, token, created_at, updated_at, ip_address, user_agent, user_id
) VALUES (
    $1, $2, NOW(), NOW(), $3, $4, $5
)
RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM sessions
WHERE token = $1 LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE token = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at < NOW();

-- =============================================================================
-- Verifications
-- =============================================================================

-- name: CreateVerification :one
INSERT INTO verifications (
    identifier, value, expires_at, created_at, updated_at
) VALUES (
    $1, $2, $3, NOW(), NOW()
)
RETURNING *;

-- name: GetVerification :one
SELECT * FROM verifications
WHERE identifier = $1 AND value = $2 LIMIT 1;

-- name: DeleteVerification :exec
DELETE FROM verifications
WHERE identifier = $1 AND value = $2;
