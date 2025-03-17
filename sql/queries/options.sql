-- name: GetOptionByID :one
SELECT id, name, created_at, updated_at, poll_id
FROM options
WHERE id = $1;

-- name: GetOptionsByPollID :many
SELECT id, name, created_at, updated_at, poll_id
FROM options
WHERE poll_id = $1;

-- name: CreateOption :one
INSERT INTO options (name, poll_id)
VALUES ($1, $2)
RETURNING id, name, created_at, updated_at, poll_id;

-- name: UpdateOption :one
UPDATE options
SET name = coalesce($2, name), updated_at = now()
WHERE id = $1
RETURNING id, name, created_at, updated_at, poll_id;

-- name: DeleteOption :exec
DELETE FROM options
WHERE id = $1;
