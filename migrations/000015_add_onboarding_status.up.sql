ALTER TABLE users ADD COLUMN onboarding_status TEXT NOT NULL DEFAULT 'NOT_STARTED';

ALTER TABLE users ADD CONSTRAINT chk_onboarding_status CHECK (
  onboarding_status IN (
    'NOT_STARTED',
    'IN_PROGRESS',
    'DOCUMENTATION_PENDING',
    'UNDER_REVIEW',
    'COMPLETED'
  )
);
