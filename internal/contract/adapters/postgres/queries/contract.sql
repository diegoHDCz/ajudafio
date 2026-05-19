-- name: CreateContract :one
INSERT INTO contracts (
  client_id,
  professional_id,
  hour_rate,
  total_amount,
  details,
  week_days,
  shift,
  start_time,
  hours_per_day,
  total_hours
) VALUES (
  @client_id,
  @professional_id,
  @hour_rate,
  @total_amount,
  @details,
  @week_days,
  @shift,
  @start_time,
  @hours_per_day,
  @total_hours
)
RETURNING id, client_id, professional_id, status, hour_rate, total_amount, details, week_days, shift, start_time, hours_per_day, total_hours, created_at;

-- name: GetContractByID :one
SELECT id, client_id, professional_id, status, hour_rate, total_amount, details, week_days, shift, start_time, hours_per_day, total_hours, created_at
FROM contracts
WHERE id = @id
LIMIT 1;

-- name: UpdateContract :one
UPDATE contracts SET
  status       = COALESCE(sqlc.narg('status'), status),
  hour_rate    = COALESCE(sqlc.narg('hour_rate'), hour_rate),
  total_amount = COALESCE(sqlc.narg('total_amount'), total_amount),
  details      = COALESCE(sqlc.narg('details'), details),
  week_days    = COALESCE(sqlc.narg('week_days'), week_days),
  shift        = COALESCE(sqlc.narg('shift'), shift),
  start_time   = COALESCE(sqlc.narg('start_time'), start_time),
  hours_per_day = COALESCE(sqlc.narg('hours_per_day'), hours_per_day),
  total_hours  = COALESCE(sqlc.narg('total_hours'), total_hours)
WHERE id = @id
RETURNING id, client_id, professional_id, status, hour_rate, total_amount, details, week_days, shift, start_time, hours_per_day, total_hours, created_at;

-- name: DeleteContract :exec
DELETE FROM contracts
WHERE id = @id;

-- name: GetContractsByUserID :many
SELECT id, client_id, professional_id, status, hour_rate, total_amount, details, week_days, shift, start_time, hours_per_day, total_hours, created_at
FROM contracts
WHERE client_id = @client_id
ORDER BY created_at DESC;

-- name: GetAllContracts :many
SELECT id, client_id, professional_id, status, hour_rate, total_amount, details, week_days, shift, start_time, hours_per_day, total_hours, created_at
FROM contracts
ORDER BY created_at DESC;

-- name: GetAllContractsByProfessionalID :many
SELECT id, client_id, professional_id, status, hour_rate, total_amount, details, week_days, shift, start_time, hours_per_day, total_hours, created_at
FROM contracts
WHERE professional_id = @professional_id
ORDER BY created_at DESC;

-- name: GetAllContractsByStatus :many
SELECT id, client_id, professional_id, status, hour_rate, total_amount, details, week_days, shift, start_time, hours_per_day, total_hours, created_at
FROM contracts
WHERE status = @status
ORDER BY created_at DESC;

-- name: GetAllContractsByUserIDAndStatus :many
SELECT id, client_id, professional_id, status, hour_rate, total_amount, details, week_days, shift, start_time, hours_per_day, total_hours, created_at
FROM contracts
WHERE client_id = @client_id
  AND status = @status
ORDER BY created_at DESC;

-- name: GetAllContractsByProfessionalIDAndStatus :many
SELECT id, client_id, professional_id, status, hour_rate, total_amount, details, week_days, shift, start_time, hours_per_day, total_hours, created_at
FROM contracts
WHERE professional_id = @professional_id
  AND status = @status
ORDER BY created_at DESC;

-- name: GetAllContractsByProfessionalCategory :many
SELECT c.id, c.client_id, c.professional_id, c.status, c.hour_rate, c.total_amount, c.details, c.week_days, c.shift, c.start_time, c.hours_per_day, c.total_hours, c.created_at
FROM contracts c
JOIN professionals p ON c.professional_id = p.id
WHERE p.category = @category
ORDER BY c.created_at DESC;
