CREATE TABLE user_roles (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role       TEXT NOT NULL CHECK (role IN ('FAMILY_CLIENT', 'HEALTH_CAREPROVIDER', 'FINANCIAL_SPONSOR', 'PLATFORM_ADMIN')),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),

  UNIQUE (user_id, role)
);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);

-- Backfill: one row per existing user, mirroring their current single `users.role`.
INSERT INTO user_roles (user_id, role)
SELECT id, role FROM users;
