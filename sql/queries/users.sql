-- name: GetUsers :many
Select
    *
FROM
    users;

-- name: GetUserById :one
SELECT
    *
FROM
    users
WHERE
    id = $1;

-- name: UpdateUser :one
UPDATE
    users
SET
    email = COALESCE($1, email),
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    hashed_password = COALESCE($4, hashed_password),
    provider = COALESCE($5, provider),
    provider_id = COALESCE($6, provider_id),
    refresh_token = COALESCE($7, refresh_token),
    role = COALESCE($8, role),
    updated_at = NOW()
WHERE id = $9 RETURNING *;


-- name: CreateUser :one
INSERT INTO
    users (email, first_name, last_name, hash_password,provider,provider_id,refresh_token,role)
VALUES
    ($1, $2, $3, $4, $5,$6,$7,$8)

RETURNING
    *;

-- name: DeleteUser :exec
DELETE FROM
    users
WHERE
    id = $1;
