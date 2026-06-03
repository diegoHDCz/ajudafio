CREATE TABLE booking_requests (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  client_id        UUID NOT NULL REFERENCES users(id),
  professional_id  UUID NOT NULL REFERENCES professionals(id),
  address_id       UUID NOT NULL REFERENCES addresses(id),
  proposed_value   NUMERIC(10,2) NOT NULL,
  schedule_details JSONB NOT NULL,
  status           VARCHAR(20) NOT NULL DEFAULT 'PENDING',
  rejection_reason TEXT,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  responded_at     TIMESTAMPTZ,

  CONSTRAINT chk_booking_status CHECK (
    status IN ('PENDING','ACCEPTED','REJECTED','EXPIRED','CANCELLED')
  ),
  CONSTRAINT chk_rejection_reason CHECK (
    status <> 'REJECTED' OR rejection_reason IS NOT NULL
  )
);

CREATE INDEX idx_booking_requests_client_id       ON booking_requests (client_id);
CREATE INDEX idx_booking_requests_professional_id ON booking_requests (professional_id);
CREATE INDEX idx_booking_requests_status          ON booking_requests (status);
