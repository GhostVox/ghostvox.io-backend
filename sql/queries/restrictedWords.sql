-- name: GetAllRestrictedWords :many
SELECT word FROM restrictedWords;

-- name: AddRestrictedWord :one
INSERT INTO restrictedWords (word) 
VALUES ($1) 
ON CONFLICT (word) DO NOTHING
RETURNING *;

-- name: AddRestrictedWordsBatch :exec
INSERT INTO restrictedWords (word)
SELECT unnest($1::text[])
ON CONFLICT (word) DO NOTHING;
