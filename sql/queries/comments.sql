-- name: GetTotalComments :one
SELECT COUNT(*) FROM comments WHERE poll_id = $1;

-- name: GetTotalCommentsByPollIDs :many
-- used by pollhandler.processPollData
SELECT poll_id, COUNT(*) as count
FROM comments
WHERE poll_id = ANY($1::uuid[])
GROUP BY poll_id;

-- name: GetAllCommentsByPollID :one
 -- in Use in commenthandler.GetAllPollComments
SELECT comments.*, users.user_name
FROM comments
JOIN users ON comments.user_id = users.id
WHERE poll_id = $1;

-- name: CreateComment :one
-- in Use in commenthandler.CreateComment
INSERT INTO comments (poll_id, user_id, content)
VALUES ($1, $2, $3)
RETURNING id;
