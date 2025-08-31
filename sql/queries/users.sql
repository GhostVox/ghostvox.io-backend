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

-- name: UpdateUserProfile :one
UPDATE
    users
SET
    email = COALESCE($1, email),
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    user_name = COALESCE($4, user_name),
    updated_at = NOW()
WHERE id = $5 RETURNING *;


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


-- name: UpdateUserAvatar :one
UPDATE users
SET picture_url = $2,
    updated_at = NOW()
where id = $1
RETURNING *;
