ALTER TABLE users ADD COLUMN auth_user_id UUID UNIQUE;

CREATE INDEX idx_users_auth_user_id ON users(auth_user_id);
