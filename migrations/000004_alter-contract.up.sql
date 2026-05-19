ALTER TABLE contracts
  ADD COLUMN week_days    VARCHAR(20)[] NOT NULL DEFAULT '{}',
  ADD COLUMN shift        VARCHAR(20),
  ADD COLUMN start_time   TIME          NOT NULL DEFAULT '08:00:00',
  ADD COLUMN hours_per_day INTEGER      NOT NULL DEFAULT 0,
  ADD COLUMN total_hours  INTEGER       NOT NULL DEFAULT 0;

ALTER TABLE contracts
  ALTER COLUMN week_days   DROP DEFAULT,
  ALTER COLUMN start_time  DROP DEFAULT,
  ALTER COLUMN hours_per_day DROP DEFAULT,
  ALTER COLUMN total_hours DROP DEFAULT;

ALTER TABLE contracts
  ADD CONSTRAINT chk_contract_shift
  CHECK (shift IS NULL OR shift IN ('MORNING','AFTERNOON','NIGHT','FULL_DAY','CUSTOM'));
