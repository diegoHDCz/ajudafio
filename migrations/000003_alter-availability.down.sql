-- 1. Remove constraints atualizadas
ALTER TABLE availabilities DROP CONSTRAINT IF EXISTS chk_day_of_week;
ALTER TABLE availabilities DROP CONSTRAINT IF EXISTS chk_shift;

-- 2. Converte day_of_week de ARRAY para VARCHAR simples (pega o primeiro elemento)
ALTER TABLE availabilities
  ALTER COLUMN day_of_week TYPE VARCHAR(10)
  USING day_of_week[1]::VARCHAR(10);

-- 3. Converte shift de ARRAY para VARCHAR simples (pega o primeiro elemento)
ALTER TABLE availabilities
  ALTER COLUMN shift TYPE VARCHAR(10)
  USING CASE WHEN shift IS NOT NULL THEN shift[1]::VARCHAR(10) ELSE NULL END;

-- 4. Restaura constraints originais
ALTER TABLE availabilities
  ADD CONSTRAINT chk_day_of_week
  CHECK (day_of_week IN ('MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY','SATURDAY','SUNDAY'));

ALTER TABLE availabilities
  ADD CONSTRAINT chk_shift
  CHECK (shift IS NULL OR shift IN ('MORNING','AFTERNOON','NIGHT','FULL_DAY'));
