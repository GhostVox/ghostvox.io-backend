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
    email = $1,
    first_name = $2,
    last_name = $3,
    user_token = $4,
    role = $5,
    updated_at = NOW()
WHERE
    id = $6 RETURNING *;

-- name: CreateUser :one
INSERT INTO
    users (email, first_name, last_name, user_token, role)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: DeleteUser :exec
DELETE FROM
    users
WHERE
    id = $1;
