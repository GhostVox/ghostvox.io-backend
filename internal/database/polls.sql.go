// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: polls.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPoll = `-- name: CreatePoll :one
INSERT INTO
    polls (user_id, title, description, expires_at, status)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING
    id, user_id, title, description, created_at, updated_at, expires_at, status
`

type CreatePollParams struct {
	UserID      string
	Title       string
	Description string
	ExpiresAt   time.Time
	Status      PollStatus
}

func (q *Queries) CreatePoll(ctx context.Context, arg CreatePollParams) (Poll, error) {
	row := q.db.QueryRowContext(ctx, createPoll,
		arg.UserID,
		arg.Title,
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
    id = $1 RETURNING id, user_id, title, description, created_at, updated_at, expires_at, status
`

func (q *Queries) DeletePoll(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePoll, id)
	return err
}

const getAllPolls = `-- name: GetAllPolls :many
SELECT
    id, user_id, title, description, created_at, updated_at, expires_at, status
FROM
    polls
`

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

const getPoll = `-- name: GetPoll :one
SELECT
    id, user_id, title, description, created_at, updated_at, expires_at, status
FROM
    polls
WHERE
    id = $1
`

func (q *Queries) GetPoll(ctx context.Context, id uuid.UUID) (Poll, error) {
	row := q.db.QueryRowContext(ctx, getPoll, id)
	var i Poll
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.Status,
	)
	return i, err
}

const getPollsByStatus = `-- name: GetPollsByStatus :many
SELECT
    id, user_id, title, description, created_at, updated_at, expires_at, status
FROM
    polls
WHERE
    status = $1
`

func (q *Queries) GetPollsByStatus(ctx context.Context, status PollStatus) ([]Poll, error) {
	rows, err := q.db.QueryContext(ctx, getPollsByStatus, status)
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
    id, user_id, title, description, created_at, updated_at, expires_at, status
FROM
    polls
WHERE
    user_id = $1
`

func (q *Queries) GetPollsByUser(ctx context.Context, userID string) ([]Poll, error) {
	rows, err := q.db.QueryContext(ctx, getPollsByUser, userID)
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

const updatePoll = `-- name: UpdatePoll :one
UPDATE
    polls
SET
    user_id = coalesce($1, user_id),
    title = coalesce($2, title),
    description = coalesce($3, description),
    expires_at = coalesce($4, expires_at),
    status = coalesce($5, status),
    updated_at = now()
WHERE
    id = $6 RETURNING id, user_id, title, description, created_at, updated_at, expires_at, status
`

type UpdatePollParams struct {
	UserID      string
	Title       string
	Description string
	ExpiresAt   time.Time
	Status      PollStatus
	ID          uuid.UUID
}

func (q *Queries) UpdatePoll(ctx context.Context, arg UpdatePollParams) (Poll, error) {
	row := q.db.QueryRowContext(ctx, updatePoll,
		arg.UserID,
		arg.Title,
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
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.Status,
	)
	return i, err
}
