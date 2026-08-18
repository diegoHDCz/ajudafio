-- name: GetUserAuthByEmail :one
SELECT id, name, email, role, password_hash
FROM users
WHERE email = @email
LIMIT 1;

-- name: SetUserPasswordHash :exec
UPDATE users SET
  password_hash = @password_hash,
  updated_at    = NOW()
WHERE id = @id;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
  id,
  user_id,
  token_hash,
  expires_at
) VALUES (
  @id,
  @user_id,
  @token_hash,
  @expires_at
)
RETURNING id, user_id, token_hash, expires_at, revoked_at, created_at;

-- name: GetRefreshTokenByHash :one
SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
FROM refresh_tokens
WHERE token_hash = @token_hash
LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET
  revoked_at = NOW()
WHERE token_hash = @token_hash;

-- name: RevokeAllRefreshTokensForUser :exec
UPDATE refresh_tokens SET
  revoked_at = NOW()
WHERE user_id = @user_id
  AND revoked_at IS NULL;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW() - INTERVAL '24 hours';

-- name: GetIdentityByProvider :one
SELECT id, user_id, provider, provider_user_id, email, created_at
FROM identities
WHERE provider = @provider
  AND provider_user_id = @provider_user_id
LIMIT 1;

-- name: CreateIdentity :one
INSERT INTO identities (
  id,
  user_id,
  provider,
  provider_user_id,
  email
) VALUES (
  @id,
  @user_id,
  @provider,
  @provider_user_id,
  @email
)
RETURNING id, user_id, provider, provider_user_id, email, created_at;
