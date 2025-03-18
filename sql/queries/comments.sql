-- name: GetTotalComments :one
SELECT COUNT(*) FROM comments WHERE poll_id = $1;
