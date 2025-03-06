-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: UpdateRefreshToken :one
UPDATE refresh_tokens SET token = $2, expires_at = $3 WHERE user_id = $1 RETURNING *;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens WHERE token = $1;


-- name: GetRefreshTokenByUserID :one
SELECT * FROM refresh_tokens WHERE user_id = $1;

-- name: DeleteRefreshTokenByUserID :exec
DELETE FROM refresh_tokens WHERE user_id = $1;
