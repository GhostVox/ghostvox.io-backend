-- name: CreatePoll :one
-- used by transactions createPollWithOptions
INSERT INTO
    polls (user_id, title, category, description, expires_at, status)
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetPollsByUser :many
-- used by pollhandler.GetPollsByUser
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
    users.last_name as CreatorLastName,
    COUNT(DISTINCT votes.id) as votes,
    COUNT(DISTINCT comments.id) as comments,
    (SELECT json_agg(options.*) FROM options WHERE options.poll_id = polls.id) as Options,
    (SELECT option_id FROM votes WHERE votes.poll_id = polls.id AND votes.user_id = $1 LIMIT 1) as UserVote
FROM
    polls
JOIN users ON polls.user_id = users.id
LEFT JOIN votes ON polls.id = votes.poll_id
LEFT JOIN comments ON polls.id = comments.poll_id
WHERE
    polls.user_id = $1 AND polls.category LIKE $2
GROUP BY
    polls.id,
    users.first_name,
    users.last_name
LIMIT $3 OFFSET $4;


-- name: UpdatePollStatus :one
-- used by cron
UPDATE
    polls
SET
    status = $2,
    updated_at = now()
WHERE
    id = $1 RETURNING *;

-- name: GetExpiredPollsToUpdate :many
-- used by cron
Select * from polls where expires_at < now() and status = 'Active';

-- name: GetAllPollsByStatusList :many
-- used by pollhandler.GetAllfinishedpolls and pollhandler.GetAllActivePolls
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
    users.last_name as CreatorLastName,
    COUNT(DISTINCT votes.id) as votes,
    COUNT(DISTINCT comments.id) as comments,
    (SELECT json_agg(options.*) FROM options WHERE options.poll_id = polls.id) as Options,
    (SELECT votes.option_id FROM votes WHERE votes.poll_id = polls.id AND votes.user_id = $5 LIMIT 1) as UserVote
FROM
    polls
JOIN users ON polls.user_id = users.id
LEFT JOIN votes ON polls.id = votes.poll_id
LEFT JOIN comments ON polls.id = comments.poll_id
WHERE
    polls.status = $1 AND polls.category LIKE($2)
GROUP BY
    polls.id,
    users.id,
    users.first_name,
    users.last_name
ORDER BY polls.expires_at DESC
LIMIT $3 OFFSET $4;

-- name: GetPollByID :one
-- name: GetPollByID :one
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
  users.last_name as CreatorLastName,
  COUNT(DISTINCT votes.id) as votes,
  COUNT(DISTINCT comments.id) as comments,
  (SELECT json_agg(options.*) FROM options WHERE options.poll_id = polls.id) as Options,
  (SELECT votes.option_id FROM votes WHERE votes.poll_id = polls.id AND votes.user_id = $2 LIMIT 1) as UserVote
FROM
  polls
  LEFT JOIN users ON polls.user_id = users.id
  LEFT JOIN votes ON polls.id = votes.poll_id
  LEFT JOIN comments ON polls.id = comments.poll_id
WHERE
  polls.id = $1
GROUP BY
  polls.id,
  users.id,
  users.first_name,
  users.last_name;

-- name: GetRecentPolls :many
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
    users.last_name as CreatorLastName,
    count(distinct votes.id) as votes,
    count(distinct comments.id) as comments,
    (SELECT json_agg(options.*) FROM options WHERE options.poll_id = polls.id) as Options,
     (SELECT votes.option_id FROM votes WHERE votes.poll_id = polls.id AND votes.user_id = $1 LIMIT 1) as UserVote
FROM polls
JOIN users ON polls.user_id = users.id
LEFT JOIN votes ON polls.id = votes.poll_id
LEFT JOIN comments ON polls.id = comments.poll_id
GROUP BY polls.id, users.first_name, users.last_name
ORDER BY polls.expires_at DESC
LIMIT 10;

--not used yet
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
