-- name: GetUserByID :one
SELECT id, auth_user_id, name, email, phone, role, onboarding_status, avatar_url, created_at, updated_at
FROM users
WHERE id = @id
LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, auth_user_id, name, email, phone, role, onboarding_status, avatar_url, created_at, updated_at
FROM users
WHERE email = @email
LIMIT 1;

-- name: GetUserByAuthUserID :one
SELECT id, auth_user_id, name, email, phone, role, onboarding_status, avatar_url, created_at, updated_at
FROM users
WHERE auth_user_id = @auth_user_id
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  id,
  auth_user_id,
  name,
  email,
  phone,
  role
) VALUES (
  @id,
  @auth_user_id,
  @name,
  @email,
  @phone,
  @role
)
RETURNING id, auth_user_id, name, email, phone, role, onboarding_status, avatar_url, created_at, updated_at;

-- name: UpdateUser :one
UPDATE users SET
  name       = COALESCE(@name, name),
  email      = COALESCE(@email, email),
  phone      = COALESCE(@phone, phone),
  role       = COALESCE(@role, role),
  updated_at = NOW()
WHERE id = @id
RETURNING id, auth_user_id, name, email, phone, role, onboarding_status, avatar_url, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = @id;

-- name: UpdateUserRole :exec
UPDATE users SET
  role       = @role,
  updated_at = NOW()
WHERE id = @id;

-- name: UpdateUserAvatar :one
UPDATE users SET
  avatar_url = @avatar_url,
  updated_at = NOW()
WHERE id = @id
RETURNING id, auth_user_id, name, email, phone, role, onboarding_status, avatar_url, created_at, updated_at;

-- name: UpdateOnboardingStatus :exec
UPDATE users SET
  onboarding_status = @onboarding_status,
  updated_at        = NOW()
WHERE id = @id;

-- name: ListRolesByUserID :many
SELECT id, user_id, role, created_at
FROM user_roles
WHERE user_id = @user_id
ORDER BY created_at;

-- name: AddUserRole :one
INSERT INTO user_roles (user_id, role)
VALUES (@user_id, @role)
ON CONFLICT (user_id, role) DO NOTHING
RETURNING id, user_id, role, created_at;

-- name: FindOrCreateUserByAuthUserID :one
-- Links a Supabase auth_user_id to an existing app user matched by email, or
-- creates a fresh one. Race-safe: two concurrent first-access requests for the
-- same brand-new email both converge on the same row via ON CONFLICT (email).
INSERT INTO users (id, auth_user_id, name, email, phone, role)
VALUES (@id, @auth_user_id, @name, @email, NULL, @role)
ON CONFLICT (email) DO UPDATE SET
  auth_user_id = COALESCE(users.auth_user_id, EXCLUDED.auth_user_id),
  updated_at   = NOW()
RETURNING id, auth_user_id, name, email, phone, role, onboarding_status, avatar_url, created_at, updated_at;
