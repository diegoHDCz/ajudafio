-- name: CreateBookingRequest :one
INSERT INTO booking_requests (
  client_id, professional_id, address_id, proposed_value, schedule_details
) VALUES (
  @client_id, @professional_id, @address_id, @proposed_value, @schedule_details
)
RETURNING id, client_id, professional_id, address_id, proposed_value, schedule_details, status, rejection_reason, created_at, responded_at;

-- name: GetBookingRequestByID :one
SELECT id, client_id, professional_id, address_id, proposed_value, schedule_details, status, rejection_reason, created_at, responded_at
FROM booking_requests
WHERE id = @id;

-- name: GetBookingRequestsByClientID :many
SELECT id, client_id, professional_id, address_id, proposed_value, schedule_details, status, rejection_reason, created_at, responded_at
FROM booking_requests
WHERE client_id = @client_id
ORDER BY created_at DESC;

-- name: GetBookingRequestsByProfessionalID :many
SELECT id, client_id, professional_id, address_id, proposed_value, schedule_details, status, rejection_reason, created_at, responded_at
FROM booking_requests
WHERE professional_id = @professional_id
ORDER BY created_at DESC;

-- name: UpdateBookingRequestStatus :one
UPDATE booking_requests
SET status           = @status,
    rejection_reason = @rejection_reason,
    responded_at     = NOW()
WHERE id = @id
RETURNING id, client_id, professional_id, address_id, proposed_value, schedule_details, status, rejection_reason, created_at, responded_at;
