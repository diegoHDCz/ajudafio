ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_onboarding_status;
ALTER TABLE users DROP COLUMN IF EXISTS onboarding_status;
