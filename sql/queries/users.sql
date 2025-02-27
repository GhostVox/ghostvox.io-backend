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

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1;

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
    role = COALESCE($7, role),
    picture_url = COALESCE($9, avatar_url),
    updated_at = NOW()
WHERE id = $8 RETURNING *;


-- name: CreateUser :one
INSERT INTO
    users (id,email, first_name, last_name, hashed_password,provider,provider_id,role,picture_url)
VALUES
<<<<<<< HEAD
    ($1, $2, $3, $4, $5,$6,$7,$8)

=======
    ($1, $2, $3, $4, $5,$6,$7,$8,$9)
>>>>>>> 54a1676 (added a few helper functions for setting cookies and adding generating and adding refresh tokens to my database.)
RETURNING
    *;

-- name: DeleteUser :exec
DELETE FROM
    users
WHERE
    id = $1;

-- name: GetUserByProviderAndProviderId :one
SELECT
    *
FROM
    users
WHERE
    provider = $1 AND provider_id = $2;
