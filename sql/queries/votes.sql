-- name: CreateVote :one
INSERT INTO votes (poll_id, option_id, user_id)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetVotesByPollID :many
SELECT * FROM votes WHERE poll_id = $1;

-- name: GetVotesByOptionID :many
SELECT * FROM votes WHERE option_id = $1;

-- name: GetVotesByUserID :many
SELECT * FROM votes WHERE user_id = $1;

-- name: DeleteVoteByID :exec
DELETE FROM votes WHERE id = $1 RETURNING *;

-- name: DeleteVotesByPollID :exec
DELETE FROM votes WHERE poll_id = $1 RETURNING *;

-- name: DeleteVotesByOptionID :exec
DELETE FROM votes WHERE option_id = $1 RETURNING *;
