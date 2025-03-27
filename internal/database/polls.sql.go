// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: polls.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPoll = `-- name: CreatePoll :one
INSERT INTO
    polls (user_id, title, category, description, expires_at, status)
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING
    id, user_id, title, description, category, created_at, updated_at, expires_at, status
`

type CreatePollParams struct {
	UserID      uuid.UUID
	Title       string
	Category    string
	Description string
	ExpiresAt   time.Time
	Status      PollStatus
}

// used by transactions createPollWithOptions
func (q *Queries) CreatePoll(ctx context.Context, arg CreatePollParams) (Poll, error) {
	row := q.db.QueryRowContext(ctx, createPoll,
		arg.UserID,
		arg.Title,
		arg.Category,
		arg.Description,
		arg.ExpiresAt,
		arg.Status,
	)
	var i Poll
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.Status,
	)
	return i, err
}

const deletePoll = `-- name: DeletePoll :exec
DELETE FROM
    polls
WHERE
    id = $1 RETURNING id, user_id, title, description, category, created_at, updated_at, expires_at, status
`

func (q *Queries) DeletePoll(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePoll, id)
	return err
}

const getAllPolls = `-- name: GetAllPolls :many
SELECT
    id, user_id, title, description, category, created_at, updated_at, expires_at, status
FROM
    polls
`

// not used yet
func (q *Queries) GetAllPolls(ctx context.Context) ([]Poll, error) {
	rows, err := q.db.QueryContext(ctx, getAllPolls)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Poll
	for rows.Next() {
		var i Poll
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Description,
			&i.Category,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExpiresAt,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllPollsByStatusList = `-- name: GetAllPollsByStatusList :many
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
    polls.status = $1 and polls.category  like($2)
    Group by polls.id, users.id
    Order by polls.expires_at desc

    limit $3 offset $4
`

type GetAllPollsByStatusListParams struct {
	Status   PollStatus
	Category string
	Limit    int32
	Offset   int32
}

type GetAllPollsByStatusListRow struct {
	Pollid           uuid.UUID
	Title            string
	Category         string
	Description      string
	Expiresat        time.Time
	Status           PollStatus
	Createdat        time.Time
	Updatedat        time.Time
	Creatorfirstname string
	Creatorlastname  sql.NullString
}

// used by pollhandler.GetAllfinishedpolls and pollhandler.GetAllActivePolls
func (q *Queries) GetAllPollsByStatusList(ctx context.Context, arg GetAllPollsByStatusListParams) ([]GetAllPollsByStatusListRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllPollsByStatusList,
		arg.Status,
		arg.Category,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllPollsByStatusListRow
	for rows.Next() {
		var i GetAllPollsByStatusListRow
		if err := rows.Scan(
			&i.Pollid,
			&i.Title,
			&i.Category,
			&i.Description,
			&i.Expiresat,
			&i.Status,
			&i.Createdat,
			&i.Updatedat,
			&i.Creatorfirstname,
			&i.Creatorlastname,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getExpiredPollsToUpdate = `-- name: GetExpiredPollsToUpdate :many
Select id, user_id, title, description, category, created_at, updated_at, expires_at, status from polls where expires_at < now() and status = 'Active'
`

// used by cron
func (q *Queries) GetExpiredPollsToUpdate(ctx context.Context) ([]Poll, error) {
	rows, err := q.db.QueryContext(ctx, getExpiredPollsToUpdate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Poll
	for rows.Next() {
		var i Poll
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Description,
			&i.Category,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExpiresAt,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPollsByUser = `-- name: GetPollsByUser :many
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
    user_id = $1 and polls.category like($2)
    limit $3 offset $4
`

type GetPollsByUserParams struct {
	UserID   uuid.UUID
	Category string
	Limit    int32
	Offset   int32
}

type GetPollsByUserRow struct {
	Pollid           uuid.UUID
	Title            string
	Category         string
	Description      string
	Expiresat        time.Time
	Status           PollStatus
	Createdat        time.Time
	Updatedat        time.Time
	Creatorfirstname string
	Creatorlastname  sql.NullString
}

// used by pollhandler.GetPollsByUser
func (q *Queries) GetPollsByUser(ctx context.Context, arg GetPollsByUserParams) ([]GetPollsByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getPollsByUser,
		arg.UserID,
		arg.Category,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPollsByUserRow
	for rows.Next() {
		var i GetPollsByUserRow
		if err := rows.Scan(
			&i.Pollid,
			&i.Title,
			&i.Category,
			&i.Description,
			&i.Expiresat,
			&i.Status,
			&i.Createdat,
			&i.Updatedat,
			&i.Creatorfirstname,
			&i.Creatorlastname,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePoll = `-- name: UpdatePoll :one
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
    id = $7 RETURNING id, user_id, title, description, category, created_at, updated_at, expires_at, status
`

type UpdatePollParams struct {
	UserID      uuid.UUID
	Title       string
	Category    string
	Description string
	ExpiresAt   time.Time
	Status      PollStatus
	ID          uuid.UUID
}

func (q *Queries) UpdatePoll(ctx context.Context, arg UpdatePollParams) (Poll, error) {
	row := q.db.QueryRowContext(ctx, updatePoll,
		arg.UserID,
		arg.Title,
		arg.Category,
		arg.Description,
		arg.ExpiresAt,
		arg.Status,
		arg.ID,
	)
	var i Poll
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.Status,
	)
	return i, err
}

const updatePollStatus = `-- name: UpdatePollStatus :one
UPDATE
    polls
SET
    status = $2,
    updated_at = now()
WHERE
    id = $1 RETURNING id, user_id, title, description, category, created_at, updated_at, expires_at, status
`

type UpdatePollStatusParams struct {
	ID     uuid.UUID
	Status PollStatus
}

// used by cron
func (q *Queries) UpdatePollStatus(ctx context.Context, arg UpdatePollStatusParams) (Poll, error) {
	row := q.db.QueryRowContext(ctx, updatePollStatus, arg.ID, arg.Status)
	var i Poll
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.Status,
	)
	return i, err
}
