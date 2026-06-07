-- name: GetProfessionalByID :one
SELECT p.id, p.user_id, p.license_number, p.category, p.years_of_experience, p.verified, p.resume, p.metadata, p.created_at, p.updated_at, u.name AS user_name, u.avatar_url AS user_avatar_url, u.email AS user_email, u.role AS user_role
FROM professionals p
INNER JOIN users u ON p.user_id = u.id
WHERE p.id = @id
LIMIT 1;

-- name: GetProfessionalByUserID :one
SELECT id, user_id, license_number, category, years_of_experience, verified, resume, metadata, created_at, updated_at
FROM professionals
WHERE user_id = @user_id
LIMIT 1;

-- name: CreateProfessional :one
INSERT INTO professionals (
  user_id,
  license_number,
  category,
  years_of_experience,
  resume,
  metadata
) VALUES (
  @user_id,
  @license_number,
  @category,
  @years_of_experience,
  @resume,
  @metadata
)
RETURNING id, user_id, license_number, category, years_of_experience, verified, resume, metadata, created_at, updated_at;

-- name: UpdateProfessional :one
UPDATE professionals SET
  license_number      = COALESCE(sqlc.narg('license_number'), license_number),
  category            = COALESCE(sqlc.narg('category'), category),
  years_of_experience = COALESCE(sqlc.narg('years_of_experience'), years_of_experience),
  verified            = COALESCE(sqlc.narg('verified'), verified),
  resume              = COALESCE(sqlc.narg('resume'), resume),
  metadata            = COALESCE(sqlc.narg('metadata'), metadata),
  updated_at          = NOW()
WHERE id = @id
RETURNING id, user_id, license_number, category, years_of_experience, verified, resume, metadata, created_at, updated_at;

-- name: DeleteProfessional :exec
DELETE FROM professionals
WHERE id = @id;

-- name: ListProfessionals :many
SELECT
    p.id,
    p.user_id,
    p.license_number,
    p.category,
    p.years_of_experience,
    p.verified,
    p.resume,
    p.metadata,
    p.created_at,
    p.updated_at,
    u.name AS user_name,
    u.avatar_url AS user_avatar_url,
    u.email AS user_email,
    u.role AS user_role
FROM professionals p
LEFT JOIN users u ON p.user_id = u.id
WHERE
    (sqlc.narg('city')::text IS NULL OR EXISTS (
        SELECT 1 FROM addresses a WHERE a.user_id = p.user_id AND a.city = sqlc.narg('city')::text
    ))
    AND (sqlc.narg('state')::text IS NULL OR EXISTS (
        SELECT 1 FROM addresses a WHERE a.user_id = p.user_id AND a.state = sqlc.narg('state')::text
    ))
    AND (sqlc.narg('day_of_week')::text[] IS NULL OR EXISTS (
        SELECT 1 FROM availabilities av WHERE av.professional_id = p.id AND av.day_of_week = ANY(sqlc.narg('day_of_week')::text[])
    ))
    AND (sqlc.narg('shift')::text[] IS NULL OR EXISTS (
        SELECT 1 FROM availabilities av WHERE av.professional_id = p.id AND av.shift = ANY(sqlc.narg('shift')::text[])
    ))
ORDER BY p.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: CountProfessionals :one
SELECT COUNT(*)
FROM professionals p
WHERE
    (sqlc.narg('city')::text IS NULL OR EXISTS (
        SELECT 1 FROM addresses a WHERE a.user_id = p.user_id AND a.city = sqlc.narg('city')::text
    ))
    AND (sqlc.narg('state')::text IS NULL OR EXISTS (
        SELECT 1 FROM addresses a WHERE a.user_id = p.user_id AND a.state = sqlc.narg('state')::text
    ))
    AND (sqlc.narg('day_of_week')::text[] IS NULL OR EXISTS (
        SELECT 1 FROM availabilities av WHERE av.professional_id = p.id AND av.day_of_week = ANY(sqlc.narg('day_of_week')::text[])
    ))
    AND (sqlc.narg('shift')::text[] IS NULL OR EXISTS (
        SELECT 1 FROM availabilities av WHERE av.professional_id = p.id AND av.shift = ANY(sqlc.narg('shift')::text[])
    ));