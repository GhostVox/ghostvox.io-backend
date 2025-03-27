-- name: GetTotalComments :one
SELECT COUNT(*) FROM comments WHERE poll_id = $1;

-- name: GetTotalCommentsByPollIDs :many
-- used by pollhandler.processPollData
SELECT poll_id, COUNT(*) as count
FROM comments
WHERE poll_id = ANY($1::uuid[])
GROUP BY poll_id;
