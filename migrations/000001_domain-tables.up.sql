-- Users
CREATE TABLE IF NOT EXISTS users (
  id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  name       TEXT        NOT NULL,
  email      TEXT        NOT NULL UNIQUE,
  phone      TEXT,
  role       VARCHAR(20) NOT NULL DEFAULT 'CLIENT',
  created_at TIMESTAMP   NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP   NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_user_role CHECK (role IN ('CLIENT', 'PROFESSIONAL', 'ADMIN'))
);

-- Professionals
CREATE TABLE professionals (
  id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id             UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  license_number      TEXT,
  category            TEXT NOT NULL,
  years_of_experience INTEGER,
  verified            BOOLEAN NOT NULL DEFAULT FALSE,
  resume              TEXT,
  metadata            JSONB,
  created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at          TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Availabilities
CREATE TABLE availabilities (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  professional_id UUID NOT NULL REFERENCES professionals(id) ON DELETE CASCADE,
  day_of_week     VARCHAR(10) NOT NULL,
  shift           VARCHAR(10),
  start_hour      VARCHAR(5),
  end_hour        VARCHAR(5),

  CONSTRAINT chk_day_of_week CHECK (day_of_week IN ('MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY','SATURDAY','SUNDAY')),
  CONSTRAINT chk_shift       CHECK (shift IS NULL OR shift IN ('MORNING','AFTERNOON','NIGHT','FULL_DAY'))
);

-- Contracts
CREATE TABLE contracts (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  client_id       UUID NOT NULL REFERENCES users(id),
  professional_id UUID NOT NULL REFERENCES professionals(id),
  status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
  hour_rate       INTEGER NOT NULL,
  total_amount    INTEGER NOT NULL,
  details         JSONB NOT NULL,
  created_at      TIMESTAMP NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_status CHECK (status IN ('PENDING','ACTIVE','COMPLETED','CANCELLED'))
);

-- Addresses
CREATE TABLE addresses (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id      UUID REFERENCES users(id) ON DELETE CASCADE,
  contract_id  UUID REFERENCES contracts(id) ON DELETE CASCADE,
  zip_code     TEXT NOT NULL,
  address_line TEXT NOT NULL,
  number       TEXT NOT NULL,
  complement   TEXT,
  district     TEXT NOT NULL,
  city         TEXT NOT NULL,
  state        TEXT NOT NULL,
  reference    TEXT,
  created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Reviews
CREATE TABLE reviews (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  client_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  professional_id UUID NOT NULL REFERENCES professionals(id) ON DELETE CASCADE,
  contract_id     UUID REFERENCES contracts(id) ON DELETE SET NULL,
  rating          INTEGER NOT NULL,
  comment         TEXT,
  created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMP NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_rating CHECK (rating BETWEEN 1 AND 5)
);

CREATE UNIQUE INDEX unique_contract_review ON reviews (contract_id);