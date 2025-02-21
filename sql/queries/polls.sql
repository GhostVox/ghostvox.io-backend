-- name: CreatePoll :one
INSERT INTO
    polls (user_id, title, description, expires_at, status)
VALUES
    ($1, $2, $3, $4, $5)
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
    user_id = $1,
    title = $2,
    description = $3,
    expires_at = $4,
    status = $5
WHERE
    id = $6 RETURNING *;

-- name: DeletePoll :exec
DELETE FROM
    polls
WHERE
    id = $1 RETURNING *;
