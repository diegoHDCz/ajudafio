-- name: CreateFamilyProfile :one
INSERT INTO family_profiles (user_id) VALUES (@user_id)
ON CONFLICT (user_id) DO UPDATE SET updated_at = NOW()
RETURNING id, user_id, created_at, updated_at;

-- name: GetFamilyProfileByUserID :one
SELECT id, user_id, created_at, updated_at
FROM family_profiles
WHERE user_id = @user_id
LIMIT 1;

-- name: CreateFinancialProfile :one
INSERT INTO financial_profiles (user_id) VALUES (@user_id)
ON CONFLICT (user_id) DO UPDATE SET updated_at = NOW()
RETURNING id, user_id, created_at, updated_at;

-- name: GetFinancialProfileByUserID :one
SELECT id, user_id, created_at, updated_at
FROM financial_profiles
WHERE user_id = @user_id
LIMIT 1;
