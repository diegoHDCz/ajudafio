ALTER TABLE contracts DROP CONSTRAINT IF EXISTS chk_contract_shift;

ALTER TABLE contracts
  DROP COLUMN IF EXISTS week_days,
  DROP COLUMN IF EXISTS shift,
  DROP COLUMN IF EXISTS start_time,
  DROP COLUMN IF EXISTS hours_per_day,
  DROP COLUMN IF EXISTS total_hours;
