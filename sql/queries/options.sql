-- name: CreateOptions :execrows
-- in use by transaction createPollWithOptions
INSERT INTO options (poll_id, name)
VALUES ($1, UNNEST($2::text[]))
RETURNING id, name, created_at, updated_at, poll_id;

-- name: GetOptionsByPollIDs :many
-- used by pollhandler.processPollData
SELECT * FROM options
WHERE poll_id = ANY($1::uuid[]);

-- name: UpdateOptionCount :one
-- in use by transaction createVoteAndUpdateOptionCount
UPDATE options
SET count = count + 1, updated_at = now()
WHERE id = $1
RETURNING id, name, created_at, updated_at, poll_id;

-- name: DeleteOption :exec
-- used by optionHandler.deleteOption
DELETE FROM options
WHERE id = $1;
