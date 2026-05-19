-- name: GetProfessionalAvailabilityByProfessionalID :many
SELECT
    a.id,
    a.professional_id,
    a.day_of_week,
    a.start_hour, -- Corrigido de start_time para start_hour
    a.end_hour   -- Corrigido de end_time para end_hour
FROM
    availabilities a  
WHERE 
    a.professional_id = @professional_id; -- Alterado de $1 para @professional_id para manter o padrão

-- name: GetAllProfessionalAvailabilities :many
SELECT
    a.id,
    a.professional_id,
    a.day_of_week,
    a.start_hour,  
    a.end_hour
FROM
    availabilities a
INNER JOIN professionals p ON a.professional_id = p.id
WHERE
    p.is_active = true
    
    -- Filtro de DayOfWeek
    AND (cardinality(@day_of_weeks::text[]) = 0 OR a.day_of_week = ANY(@day_of_weeks::text[]))
    
    -- Filtro de Horário
    AND (sqlc.narg('start_time')::text IS NULL OR a.start_hour >= sqlc.narg('start_time')::text)
    AND (sqlc.narg('end_time')::text IS NULL OR a.end_hour <= sqlc.narg('end_time')::text)
    
    -- Filtro de Cidade
    AND (sqlc.narg('city')::text IS NULL OR EXISTS (
        SELECT 1
        FROM addresses ad
        WHERE ad.user_id = p.user_id
          AND ad.city = sqlc.narg('city')::text
    ));

-- name: CreateProfessionalAvailability :one
INSERT INTO availabilities (
    professional_id,
    day_of_week,
    start_hour,
    end_hour
) VALUES (
    @professional_id,
    @day_of_week,
    @start_hour,
    @end_hour
) RETURNING id, professional_id, day_of_week, start_hour, end_hour;

-- name: UpdateProfessionalAvailability :one
UPDATE availabilities SET
    day_of_week = COALESCE(sqlc.narg('day_of_week'), day_of_week),
    start_hour  = COALESCE(sqlc.narg('start_hour'), start_hour),
    end_hour    = COALESCE(sqlc.narg('end_hour'), end_hour)
WHERE id = @id
RETURNING id, professional_id, day_of_week, start_hour, end_hour;

-- name: DeleteProfessionalAvailability :exec
DELETE FROM availabilities
WHERE id = @id;