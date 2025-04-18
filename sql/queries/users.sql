-- name: GetUsers :many
Select
    *
FROM
    users;

-- name: GetUserById :one
SELECT
    *
FROM
    users
WHERE
    id = $1;

-- name: GetUserByEmail :one
-- used by auth handler
SELECT
    *
FROM
    users
WHERE
    email = $1;

-- name: UpdateUser :one
UPDATE
    users
SET
    email = COALESCE($1, email),
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    hashed_password = COALESCE($4, hashed_password),
    provider = COALESCE($5, provider),
    provider_id = COALESCE($6, provider_id),
    role = COALESCE($7, role),
    picture_url = COALESCE($9, avatar_url),
    updated_at = NOW()
WHERE id = $8 RETURNING *;


-- name: CreateUser :one
-- used by auth handler
INSERT INTO
    users (email, first_name, last_name, hashed_password,provider,provider_id,role,picture_url)
VALUES

    ($1, $2, $3, $4, $5,$6,$7,$8)
RETURNING
    *;

-- name: DeleteUser :exec
DELETE FROM
    users
WHERE
    id = $1;

-- name: GetUserByProviderAndProviderId :one
-- used by auth handler
SELECT
    *
FROM
    users
WHERE
    provider = $1 AND provider_id = $2;

-- name: GetUserStats :one
SELECT
    (SELECT COUNT(*) FROM polls WHERE polls.user_id = $1) as total_polls,
    (SELECT COUNT(*) FROM comments WHERE comments.poll_id IN (SELECT id FROM polls WHERE polls.user_id = $1)) as total_comments,
    (SELECT COUNT(*) FROM votes WHERE votes.poll_id IN (SELECT id FROM polls WHERE polls.user_id = $1)) as total_votes
FROM users
WHERE users.id = $1;

-- name: UpdateUserName :one
UPDATE
    users
SET
    user_name =  $1,
    updated_at = NOW()
WHERE id = $2 RETURNING *;

-- name: CheckUserNameExists :one
SELECT EXISTS(
    Select 1
FROM
    users
WHERE
    user_name = $1
    ) as exists;
