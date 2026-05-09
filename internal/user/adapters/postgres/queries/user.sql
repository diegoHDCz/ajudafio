-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = @id
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = @email
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  email,
  name,
  telephone,
  telephone_whatsapp,
  second_telephone,
  second_telephone_whatsapp,
  linkedin,
  instagram,
  facebook,
  identification_number,
  identification_type,
  role
) VALUES (
  @email,
  @name,
  @telephone,
  @telephone_whatsapp,
  @second_telephone,
  @second_telephone_whatsapp,
  @linkedin,
  @instagram,
  @facebook,
  @identification_number,
  @identification_type,
  @role
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
  name                      = COALESCE(@name, name),
  telephone                 = COALESCE(@telephone, telephone),
  telephone_whatsapp        = COALESCE(@telephone_whatsapp, telephone_whatsapp),
  second_telephone          = COALESCE(@second_telephone, second_telephone),
  second_telephone_whatsapp = COALESCE(@second_telephone_whatsapp, second_telephone_whatsapp),
  linkedin                  = COALESCE(@linkedin, linkedin),
  instagram                 = COALESCE(@instagram, instagram),
  facebook                  = COALESCE(@facebook, facebook),
  identification_number     = COALESCE(@identification_number, identification_number),
  identification_type       = COALESCE(@identification_type, identification_type),
  updated_at                = NOW()
WHERE id = @id
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = @id;