-- name: CreatePoll :one
INSERT INTO
    polls (user_id, title, category, description, expires_at, status)
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetPoll :one
SELECT
    *
FROM
    polls
WHERE
    id = $1;

-- name: GetPollsByUser :many
SELECT
polls.id as PollId,
    polls.title as Title,
    polls.category as Category,
    polls.description as Description,
    polls.expires_at as ExpiresAt,
    polls.status as Status,
    polls.created_at as CreatedAt,
    polls.updated_at as UpdatedAt,
    users.first_name as CreatorFirstName,
    users.last_name as CreatorLastName
FROM
    polls join users on polls.user_id = users.id
WHERE
    user_id = $1
    limit $2 offset $3;

-- name: GetPollsByStatus :many
SELECT
    *
FROM
    polls
WHERE
    status = $1;

-- name: GetAllPolls :many
SELECT
    *
FROM
    polls;

-- name: UpdatePoll :one
UPDATE
    polls
SET
    user_id = coalesce($1, user_id),
    title = coalesce($2, title),
    category = coalesce($3, category),
    description = coalesce($4, description),
    expires_at = coalesce($5, expires_at),
    status = coalesce($6, status),
    updated_at = now()
WHERE
    id = $7 RETURNING *;

-- name: DeletePoll :exec
DELETE FROM
    polls
WHERE
    id = $1 RETURNING *;

-- name: UpdatePollStatus :one
UPDATE
    polls
SET
    status = $2,
    updated_at = now()
WHERE
    id = $1 RETURNING *;

-- name: GetExpiredPollsToUpdate :many
Select * from polls where expires_at < now() and status = 'Active';

-- name: GetAllPollsByStatusList :many
SELECT
    polls.id as PollId,
    polls.title as Title,
    polls.category as Category,
    polls.description as Description,
    polls.expires_at as ExpiresAt,
    polls.status as Status,
    polls.created_at as CreatedAt,
    polls.updated_at as UpdatedAt,
    users.first_name as CreatorFirstName,
    users.last_name as CreatorLastName

FROM
    polls join users on polls.user_id = users.id
WHERE
    polls.status = $1
    Group by polls.id, users.id
    Order by polls.expires_at desc

    limit $2 offset $3
    ;
