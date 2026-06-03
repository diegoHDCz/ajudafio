-- name: GetAvailabilityByID :one
SELECT id, professional_id, day_of_week, shift, start_hour, end_hour
FROM availabilities
WHERE id = @id;

-- name: GetProfessionalAvailabilityByProfessionalID :many
SELECT id, professional_id, day_of_week, shift, start_hour, end_hour
FROM availabilities
WHERE professional_id = @professional_id;

-- name: GetAllProfessionalAvailabilities :many
SELECT a.id, a.professional_id, a.day_of_week, a.shift, a.start_hour, a.end_hour
FROM availabilities a
INNER JOIN professionals p ON a.professional_id = p.id
WHERE
    p.is_active = true
    AND (sqlc.narg('day_of_week')::text IS NULL OR a.day_of_week = sqlc.narg('day_of_week')::text)
    AND (sqlc.narg('shift')::text IS NULL OR a.shift = sqlc.narg('shift')::text)
    AND (sqlc.narg('city')::text IS NULL OR EXISTS (
        SELECT 1 FROM addresses ad
        WHERE ad.user_id = p.user_id AND ad.city = sqlc.narg('city')::text
    ));

-- name: CreateProfessionalAvailability :one
INSERT INTO availabilities (professional_id, day_of_week, shift, start_hour, end_hour)
VALUES (@professional_id, @day_of_week, @shift, @start_hour, @end_hour)
RETURNING id, professional_id, day_of_week, shift, start_hour, end_hour;

-- name: UpdateProfessionalAvailability :one
UPDATE availabilities SET
    day_of_week = COALESCE(sqlc.narg('day_of_week'), day_of_week),
    shift       = COALESCE(sqlc.narg('shift'), shift),
    start_hour  = COALESCE(sqlc.narg('start_hour'), start_hour),
    end_hour    = COALESCE(sqlc.narg('end_hour'), end_hour),
    updated_at  = NOW()
WHERE id = @id
RETURNING id, professional_id, day_of_week, shift, start_hour, end_hour;

-- name: DeleteProfessionalAvailability :exec
DELETE FROM availabilities WHERE id = @id;
