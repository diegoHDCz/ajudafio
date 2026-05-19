-- 1. Remove constraints antigas
ALTER TABLE availabilities DROP CONSTRAINT IF EXISTS chk_day_of_week;
ALTER TABLE availabilities DROP CONSTRAINT IF EXISTS chk_shift;

-- 2. Altera day_of_week para ARRAY de VARCHAR
ALTER TABLE availabilities
  ALTER COLUMN day_of_week TYPE VARCHAR(20)[]
  USING ARRAY[day_of_week]::VARCHAR(20)[];

-- 3. Altera shift para ARRAY de VARCHAR
ALTER TABLE availabilities
  ALTER COLUMN shift TYPE VARCHAR(20)[]
  USING CASE WHEN shift IS NOT NULL THEN ARRAY[shift]::VARCHAR(20)[] ELSE NULL END;

-- 4. Cria nova constraint para day_of_week
ALTER TABLE availabilities
  ADD CONSTRAINT chk_day_of_week
  CHECK (day_of_week <@ ARRAY['MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY','SATURDAY','SUNDAY']::VARCHAR(20)[]);

-- 5. Cria nova constraint para shift
ALTER TABLE availabilities
  ADD CONSTRAINT chk_shift
  CHECK (shift IS NULL OR shift <@ ARRAY['MORNING','AFTERNOON','NIGHT','FULL_DAY','CUSTOM']::VARCHAR(20)[]);
