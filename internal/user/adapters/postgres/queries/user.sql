-- name: GetUserByID :one
SELECT id, name, email, phone, role, created_at, updated_at
FROM users
WHERE id = @id
LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, name, email, phone, role, created_at, updated_at
FROM users
WHERE email = @email
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  name,
  email,
  phone,
  role
) VALUES (
  @name,
  @email,
  @phone,
  @role
)
RETURNING id, name, email, phone, role, created_at, updated_at;

-- name: UpdateUser :one
UPDATE users SET
  name       = COALESCE(@name, name),
  email      = COALESCE(@email, email),
  phone      = COALESCE(@phone, phone),
  role       = COALESCE(@role, role),
  updated_at = NOW()
WHERE id = @id
RETURNING id, name, email, phone, role, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = @id;