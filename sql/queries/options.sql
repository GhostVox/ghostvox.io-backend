-- name: GetOptionByID :one
SELECT id, name, created_at, updated_at, poll_id
FROM options
WHERE id = $1;

-- name: GetOptionsByPollID :many
SELECT id, name, count, created_at, updated_at, poll_id
FROM options
WHERE poll_id = $1;

-- name: CreateOptions :execrows
INSERT INTO options (poll_id, name)
VALUES ($1, UNNEST($2::text[]))
RETURNING id, name, created_at, updated_at, poll_id;

-- name: UpdateOption :one
UPDATE options
SET name = coalesce($2, name), updated_at = now()
WHERE id = $1
RETURNING id, name, created_at, updated_at, poll_id;

-- name: DeleteOption :exec
DELETE FROM options
WHERE id = $1;

-- name: GetOptionsByPollIDs :many
SELECT * FROM options
WHERE poll_id = ANY($1::uuid[]);
