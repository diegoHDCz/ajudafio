-- name: GetAddressByID :one
SELECT id, user_id, contract_id, zip_code, address_line, number, complement, district, city, state, reference, created_at, updated_at
FROM addresses
WHERE id = @id
LIMIT 1;

-- name: GetAddressesByUserID :many
SELECT id, user_id, contract_id, zip_code, address_line, number, complement, district, city, state, reference, created_at, updated_at
FROM addresses
WHERE user_id = @user_id
ORDER BY created_at DESC;

-- name: GetAddressesByContractID :many
SELECT id, user_id, contract_id, zip_code, address_line, number, complement, district, city, state, reference, created_at, updated_at
FROM addresses
WHERE contract_id = @contract_id
ORDER BY created_at DESC;

-- name: GetAllAddresses :many
SELECT id, user_id, contract_id, zip_code, address_line, number, complement, district, city, state, reference, created_at, updated_at
FROM addresses
WHERE
  (sqlc.narg('zip_code')::text    IS NULL OR zip_code    = sqlc.narg('zip_code')::text)
  AND (sqlc.narg('user_id')::uuid     IS NULL OR user_id     = sqlc.narg('user_id')::uuid)
  AND (sqlc.narg('city')::text        IS NULL OR city        = sqlc.narg('city')::text)
  AND (sqlc.narg('district')::text    IS NULL OR district    = sqlc.narg('district')::text)
  AND (sqlc.narg('number')::text      IS NULL OR number      = sqlc.narg('number')::text)
  AND (sqlc.narg('contract_id')::uuid IS NULL OR contract_id = sqlc.narg('contract_id')::uuid)
ORDER BY created_at DESC;

-- name: GetAddressesByCity :many
SELECT id, user_id, contract_id, zip_code, address_line, number, complement, district, city, state, reference, created_at, updated_at
FROM addresses
WHERE city = @city
ORDER BY created_at DESC;

-- name: CreateAddress :one
INSERT INTO addresses (
  user_id,
  contract_id,
  zip_code,
  address_line,
  number,
  complement,
  district,
  city,
  state,
  reference
) VALUES (
  @user_id,
  @contract_id,
  @zip_code,
  @address_line,
  @number,
  @complement,
  @district,
  @city,
  @state,
  @reference
)
RETURNING id, user_id, contract_id, zip_code, address_line, number, complement, district, city, state, reference, created_at, updated_at;

-- name: UpdateAddress :one
UPDATE addresses SET
  zip_code     = COALESCE(sqlc.narg('zip_code'), zip_code),
  address_line = COALESCE(sqlc.narg('address_line'), address_line),
  number       = COALESCE(sqlc.narg('number'), number),
  complement   = COALESCE(sqlc.narg('complement'), complement),
  district     = COALESCE(sqlc.narg('district'), district),
  city         = COALESCE(sqlc.narg('city'), city),
  state        = COALESCE(sqlc.narg('state'), state),
  reference    = COALESCE(sqlc.narg('reference'), reference),
  updated_at   = NOW()
WHERE id = @id
RETURNING id, user_id, contract_id, zip_code, address_line, number, complement, district, city, state, reference, created_at, updated_at;

-- name: DeleteAddress :exec
DELETE FROM addresses
WHERE id = @id;

-- name: GetAddressWithUser :one
SELECT
  a.id,
  a.user_id,
  a.contract_id,
  a.zip_code,
  a.address_line,
  a.number,
  a.complement,
  a.district,
  a.city,
  a.state,
  a.reference,
  a.created_at,
  a.updated_at,
  u.name        AS user_name,
  u.email       AS user_email,
  u.phone       AS user_phone,
  u.role        AS user_role
FROM addresses a
INNER JOIN users u ON u.id = a.user_id
WHERE a.id = @id
LIMIT 1;

-- name: GetAddressWithContract :one
SELECT
  a.id,
  a.user_id,
  a.contract_id,
  a.zip_code,
  a.address_line,
  a.number,
  a.complement,
  a.district,
  a.city,
  a.state,
  a.reference,
  a.created_at,
  a.updated_at,
  c.client_id          AS contract_client_id,
  c.professional_id    AS contract_professional_id,
  c.status             AS contract_status,
  c.hour_rate          AS contract_hour_rate,
  c.total_amount       AS contract_total_amount,
  c.details            AS contract_details,
  c.created_at         AS contract_created_at
FROM addresses a
INNER JOIN contracts c ON c.id = a.contract_id
WHERE a.id = @id
LIMIT 1;
