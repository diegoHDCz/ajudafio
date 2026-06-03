DROP TABLE availabilities;

CREATE TABLE availabilities (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    professional_id UUID NOT NULL REFERENCES professionals(id) ON DELETE CASCADE,
    day_of_week     VARCHAR(20)[] NOT NULL,
    shift           VARCHAR(20)[],
    start_hour      VARCHAR(5),
    end_hour        VARCHAR(5),

    CONSTRAINT chk_day_of_week CHECK (day_of_week <@ ARRAY['MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY','SATURDAY','SUNDAY']::VARCHAR(20)[]),
    CONSTRAINT chk_shift CHECK (shift IS NULL OR shift <@ ARRAY['MORNING','AFTERNOON','NIGHT','FULL_DAY','CUSTOM']::VARCHAR(20)[])
);
