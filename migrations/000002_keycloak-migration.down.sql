-- Link users to their Keycloak identity (JWT sub claim)
ALTER TABLE users ADD COLUMN keycloak_id TEXT UNIQUE;

-- Keycloak now owns session management
DROP TABLE IF EXISTS sessions;

-- Keycloak now owns OAuth provider accounts
DROP TABLE IF EXISTS accounts;

-- Keycloak now owns email verification
DROP TABLE IF EXISTS verifications;