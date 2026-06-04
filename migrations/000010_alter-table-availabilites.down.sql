-- 1. Garante que não existem valores nulos antes de aplicar o NOT NULL
-- (Substitua 'CUSTOM' pelo turno padrão que fizer mais sentido para o seu negócio)
UPDATE availabilities 
SET shift = 'CUSTOM' 
WHERE shift IS NULL;

-- 2. Se você criou o índice parcial para nulos sugerido anteriormente, remova-o aqui
DROP INDEX IF EXISTS uq_professional_day_null_shift;

-- 3. Remove as constraints alteradas
ALTER TABLE availabilities 
    DROP CONSTRAINT chk_shift,
    DROP CONSTRAINT chk_custom_hours,
    DROP CONSTRAINT uq_professional_day_shift;

-- 4. Restaura a tabela exatamente ao estado original
ALTER TABLE availabilities 
    ALTER COLUMN shift SET NOT NULL,
    ADD CONSTRAINT chk_shift CHECK (shift IN ('MORNING','AFTERNOON','NIGHT','FULL_DAY','CUSTOM')),
    ADD CONSTRAINT chk_custom_hours CHECK (shift != 'CUSTOM' OR (start_hour IS NOT NULL AND end_hour IS NOT NULL)),
    ADD CONSTRAINT uq_professional_day_shift UNIQUE (professional_id, day_of_week, shift);