-- Users
CREATE TABLE users (
  id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email                     TEXT NOT NULL UNIQUE,
  name                      TEXT,
  email_verified            BOOLEAN NOT NULL DEFAULT FALSE,
  image                     TEXT,
  telephone                 TEXT,
  telephone_whatsapp        BOOLEAN NOT NULL DEFAULT FALSE,
  second_telephone          TEXT,
  second_telephone_whatsapp BOOLEAN NOT NULL DEFAULT FALSE,
  linkedin                  TEXT,
  instagram                 TEXT,
  facebook                  TEXT,
  identification_number     TEXT UNIQUE,
  identification_type       VARCHAR(20),
  role                      VARCHAR(20) NOT NULL DEFAULT 'CLIENT',
  created_at                TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at                TIMESTAMP NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_role CHECK (role IN ('CLIENT', 'PROFESSIONAL', 'ADMIN'))
);

-- Sessions
CREATE TABLE sessions (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  expires_at  TIMESTAMP NOT NULL,
  token       TEXT NOT NULL UNIQUE,
  created_at  TIMESTAMP NOT NULL,
  updated_at  TIMESTAMP NOT NULL,
  ip_address  TEXT,
  user_agent  TEXT,
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- Accounts
CREATE TABLE accounts (
  id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  account_id                TEXT NOT NULL,
  provider_id               TEXT NOT NULL,
  user_id                   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  access_token              TEXT,
  refresh_token             TEXT,
  id_token                  TEXT,
  access_token_expires_at   TIMESTAMP,
  refresh_token_expires_at  TIMESTAMP,
  scope                     TEXT,
  password                  TEXT,
  created_at                TIMESTAMP NOT NULL,
  updated_at                TIMESTAMP NOT NULL
);

-- Verifications
CREATE TABLE verifications (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  identifier  TEXT NOT NULL,
  value       TEXT NOT NULL,
  expires_at  TIMESTAMP NOT NULL,
  created_at  TIMESTAMP,
  updated_at  TIMESTAMP
);