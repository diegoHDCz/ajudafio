-- name: CreateAppointment :one
INSERT INTO appointments (
  contract_id, client_id, professional_id,
  date, start_time, end_time,
  zip_code, address_line, number, complement, district, city, state, reference
) VALUES (
  @contract_id, @client_id, @professional_id,
  @date, @start_time, @end_time,
  @zip_code, @address_line, @number, @complement, @district, @city, @state, @reference
)
RETURNING id, contract_id, client_id, professional_id, date, start_time, end_time, status, zip_code, address_line, number, complement, district, city, state, reference, version, created_at, updated_at;

-- name: GetAppointmentByID :one
SELECT id, contract_id, client_id, professional_id, date, start_time, end_time, status, zip_code, address_line, number, complement, district, city, state, reference, version, created_at, updated_at
FROM appointments
WHERE id = @id;

-- name: GetAppointmentsByContractID :many
SELECT id, contract_id, client_id, professional_id, date, start_time, end_time, status, zip_code, address_line, number, complement, district, city, state, reference, version, created_at, updated_at
FROM appointments
WHERE contract_id = @contract_id
ORDER BY date, start_time;

-- name: GetAppointmentsByClientID :many
SELECT id, contract_id, client_id, professional_id, date, start_time, end_time, status, zip_code, address_line, number, complement, district, city, state, reference, version, created_at, updated_at
FROM appointments
WHERE client_id = @client_id
ORDER BY date, start_time;

-- name: GetAppointmentsByProfessionalID :many
SELECT id, contract_id, client_id, professional_id, date, start_time, end_time, status, zip_code, address_line, number, complement, district, city, state, reference, version, created_at, updated_at
FROM appointments
WHERE professional_id = @professional_id
ORDER BY date, start_time;

-- name: UpdateAppointmentStatus :one
UPDATE appointments
SET status     = @status,
    version    = version + 1,
    updated_at = NOW()
WHERE id = @id AND version = @version
RETURNING id, contract_id, client_id, professional_id, date, start_time, end_time, status, zip_code, address_line, number, complement, district, city, state, reference, version, created_at, updated_at;

-- name: DeleteAppointment :exec
DELETE FROM appointments WHERE id = @id;

-- name: CheckOverlap :one
SELECT EXISTS (
  SELECT 1 FROM appointments
  WHERE professional_id = @professional_id
    AND date = @date
    AND status <> 'CANCELLED'
    AND NOT (end_time <= @start_time OR start_time >= @end_time)
) AS has_overlap;
