-- name: GetUsers :many
Select
    *
FROM
    users;

-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    email = $1;

-- Name: UpdateUser :one
UPDATE
    users
SET
    email = $1,
    name = $2,
    password = $3,
    role = $4
WHERE
    id = $5;
