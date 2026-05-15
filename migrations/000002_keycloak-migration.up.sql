
ALTER TABLE users DROP COLUMN IF EXISTS keycloak_id;

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

CREATE TABLE verifications (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  identifier  TEXT NOT NULL,
  value       TEXT NOT NULL,
  expires_at  TIMESTAMP NOT NULL,
  created_at  TIMESTAMP,
  updated_at  TIMESTAMP
);
