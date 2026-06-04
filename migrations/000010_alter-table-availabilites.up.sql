ALTER TABLE availabilities 

    ALTER COLUMN shift DROP NOT NULL,
    DROP CONSTRAINT chk_shift,
    DROP CONSTRAINT chk_custom_hours,
    DROP CONSTRAINT uq_professional_day_shift;


ALTER TABLE availabilities
    ADD CONSTRAINT chk_shift CHECK (shift IS NULL OR shift IN ('MORNING','AFTERNOON','NIGHT','FULL_DAY','CUSTOM')),
    ADD CONSTRAINT chk_custom_hours CHECK (shift IS NULL OR shift != 'CUSTOM' OR (start_hour IS NOT NULL AND end_hour IS NOT NULL)),
    ADD CONSTRAINT uq_professional_day_shift UNIQUE (professional_id, day_of_week, shift);