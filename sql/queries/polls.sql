-- name: CreatePoll :one
INSERT INTO
    polls (user_id, title, category, description, expires_at, status)
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetPoll :one
SELECT
    *
FROM
    polls
WHERE
    id = $1;

-- name: GetPollsByUser :many
SELECT
    *
FROM
    polls
WHERE
    user_id = $1;

-- name: GetPollsByStatus :many
SELECT
    *
FROM
    polls
WHERE
    status = $1;

-- name: GetAllPolls :many
SELECT
    *
FROM
    polls;

-- name: UpdatePoll :one
UPDATE
    polls
SET
    user_id = coalesce($1, user_id),
    title = coalesce($2, title),
    category = coalesce($3, category),
    description = coalesce($4, description),
    expires_at = coalesce($5, expires_at),
    status = coalesce($6, status),
    updated_at = now()
WHERE
    id = $7 RETURNING *;

-- name: DeletePoll :exec
DELETE FROM
    polls
WHERE
    id = $1 RETURNING *;
