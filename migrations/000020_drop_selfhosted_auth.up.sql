-- Run only after users have been provisioned in Supabase Auth (see ADR-005 /
-- docs/result/refactor-auth.md open question #2). Irreversible in practice: the
-- down migration restores the columns/tables but not the bcrypt hashes or live
-- refresh tokens that existed before this ran.

DROP TABLE IF EXISTS identities;
DROP TABLE IF EXISTS refresh_tokens;
ALTER TABLE users DROP COLUMN IF EXISTS password_hash;
