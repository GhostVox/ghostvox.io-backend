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
