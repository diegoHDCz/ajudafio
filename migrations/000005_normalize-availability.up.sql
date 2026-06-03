DROP TABLE availabilities;

CREATE TABLE availabilities (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    professional_id UUID NOT NULL REFERENCES professionals(id) ON DELETE CASCADE,
    day_of_week     VARCHAR(20) NOT NULL,
    shift           VARCHAR(20) NOT NULL,
    start_hour      VARCHAR(5),
    end_hour        VARCHAR(5),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_day_of_week CHECK (day_of_week IN ('MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY','SATURDAY','SUNDAY')),
    CONSTRAINT chk_shift CHECK (shift IN ('MORNING','AFTERNOON','NIGHT','FULL_DAY','CUSTOM')),
    CONSTRAINT chk_custom_hours CHECK (shift != 'CUSTOM' OR (start_hour IS NOT NULL AND end_hour IS NOT NULL)),
    CONSTRAINT uq_professional_day_shift UNIQUE (professional_id, day_of_week, shift)
);
