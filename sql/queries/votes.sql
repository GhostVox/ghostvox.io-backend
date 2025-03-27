-- name: CreateVote :one
-- in use in transaction CreateVoteAndUpdateOptionCount
INSERT INTO votes (poll_id, option_id, user_id)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetTotalVotesByPollIDs :many
-- used by pollhandler.processPollData
SELECT poll_id, COUNT(*) as count
FROM votes
WHERE poll_id = ANY($1::uuid[])
GROUP BY poll_id;

-- name: GetTotalVotesByPollID :one
SELECT count(*) FROM votes WHERE poll_id = $1;

-- name: GetVotesByOptionID :many
SELECT * FROM votes WHERE option_id = $1;

-- name: GetVotesByUserID :many
SELECT * FROM votes WHERE user_id = $1;
