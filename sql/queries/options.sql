-- name: GetOptionByID :one
SELECT id, name, value, created_at, updated_at, poll_id
FROM options
WHERE id = $1;

-- name: GetOptionsByPollID :many
SELECT id, name, value, created_at, updated_at, poll_id
FROM options
WHERE poll_id = $1;

-- name: CreateOption :one
INSERT INTO options (name, value, poll_id)
VALUES ($1, $2, $3)
RETURNING id, name, value, created_at, updated_at, poll_id;

-- name: UpdateOption :one
UPDATE options
SET name = coalesce($2, name), value = coalesce($3, value), updated_at = now()
WHERE id = $1
RETURNING id, name, value, created_at, updated_at, poll_id;

-- name: DeleteOption :exec
DELETE FROM options
WHERE id = $1;
