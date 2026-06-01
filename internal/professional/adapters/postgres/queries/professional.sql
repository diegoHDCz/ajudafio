-- name: GetProfessionalByID :one
SELECT id, user_id, license_number, category, years_of_experience, verified, resume, metadata, created_at, updated_at
FROM professionals
WHERE id = @id
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
SELECT DISTINCT
    p.id,
    p.user_id,
    p.license_number,
    p.category,
    p.years_of_experience,
    p.verified,
    p.resume,
    p.metadata,
    p.created_at,
    p.updated_at
FROM professionals p
LEFT JOIN addresses a ON p.user_id = a.user_id
LEFT JOIN availabilities av ON p.id = av.professional_id
WHERE
    (sqlc.narg('city')::text IS NULL OR a.city = sqlc.narg('city')::text)
    AND (sqlc.narg('state')::text IS NULL OR a.state = sqlc.narg('state')::text)
    AND (sqlc.narg('day_of_week')::text[] IS NULL OR av.day_of_week && sqlc.narg('day_of_week')::text[])
    AND (sqlc.narg('shift')::text[] IS NULL OR av.shift && sqlc.narg('shift')::text[])
ORDER BY p.created_at DESC;