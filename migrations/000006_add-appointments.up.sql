CREATE TABLE appointments (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  contract_id     UUID NOT NULL REFERENCES contracts(id) ON DELETE RESTRICT,
  client_id       UUID NOT NULL REFERENCES users(id),
  professional_id UUID NOT NULL REFERENCES professionals(id),
  date            DATE NOT NULL,
  start_time      TIME NOT NULL,
  end_time        TIME NOT NULL,
  status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
  zip_code        TEXT NOT NULL,
  address_line    TEXT NOT NULL,
  number          TEXT NOT NULL,
  complement      TEXT,
  district        TEXT NOT NULL,
  city            TEXT NOT NULL,
  state           TEXT NOT NULL,
  reference       TEXT,
  version         INTEGER NOT NULL DEFAULT 1,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_appointment_status CHECK (status IN ('PENDING','CONFIRMED','CANCELLED','COMPLETED')),
  CONSTRAINT chk_time_order CHECK (end_time > start_time)
);

CREATE UNIQUE INDEX uq_professional_date_start
  ON appointments (professional_id, date, start_time)
  WHERE status <> 'CANCELLED';
