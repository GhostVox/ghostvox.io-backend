// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: options_votes.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const countVotesByOptionAndPollID = `-- name: CountVotesByOptionAndPollID :many
WITH vote_counts AS (
    SELECT
        option_id,
        COUNT(*) as vote_count
    FROM
        votes
    WHERE
        votes.poll_id = $1
    GROUP BY
        votes.option_id
)
SELECT
    option.name,
    COALESCE(vc.vote_count, 0) as vote_count
FROM
    options option
    LEFT JOIN vote_counts vc ON option.id = vc.option_id
WHERE
    option.poll_id = $2
`

type CountVotesByOptionAndPollIDParams struct {
	PollID   uuid.UUID
	PollID_2 uuid.UUID
}

type CountVotesByOptionAndPollIDRow struct {
	Name      string
	VoteCount int64
}

func (q *Queries) CountVotesByOptionAndPollID(ctx context.Context, arg CountVotesByOptionAndPollIDParams) ([]CountVotesByOptionAndPollIDRow, error) {
	rows, err := q.db.QueryContext(ctx, countVotesByOptionAndPollID, arg.PollID, arg.PollID_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CountVotesByOptionAndPollIDRow
	for rows.Next() {
		var i CountVotesByOptionAndPollIDRow
		if err := rows.Scan(&i.Name, &i.VoteCount); err != nil {
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

const getVotesByOptionAndPollID = `-- name: GetVotesByOptionAndPollID :many
SELECT
    votes.id, votes.poll_id, votes.option_id, votes.created_at, votes.user_id,
    op.id, op.name, op.poll_id, op.count, op.created_at, op.updated_at,
    po.id, po.user_id, po.title, po.description, po.category, po.created_at, po.updated_at, po.expires_at, po.status
FROM
    votes
    JOIN options op ON votes.option_id = op.id
    JOIN polls po ON votes.poll_id = po.id
WHERE
    votes.option_id = $1
    AND votes.poll_id = $2
`

type GetVotesByOptionAndPollIDParams struct {
	OptionID uuid.UUID
	PollID   uuid.UUID
}

type GetVotesByOptionAndPollIDRow struct {
	ID          uuid.UUID
	PollID      uuid.UUID
	OptionID    uuid.UUID
	CreatedAt   time.Time
	UserID      uuid.UUID
	ID_2        uuid.UUID
	Name        string
	PollID_2    uuid.UUID
	Count       int32
	CreatedAt_2 time.Time
	UpdatedAt   time.Time
	ID_3        uuid.UUID
	UserID_2    uuid.UUID
	Title       string
	Description string
	Category    string
	CreatedAt_3 time.Time
	UpdatedAt_2 time.Time
	ExpiresAt   time.Time
	Status      PollStatus
}

func (q *Queries) GetVotesByOptionAndPollID(ctx context.Context, arg GetVotesByOptionAndPollIDParams) ([]GetVotesByOptionAndPollIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getVotesByOptionAndPollID, arg.OptionID, arg.PollID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetVotesByOptionAndPollIDRow
	for rows.Next() {
		var i GetVotesByOptionAndPollIDRow
		if err := rows.Scan(
			&i.ID,
			&i.PollID,
			&i.OptionID,
			&i.CreatedAt,
			&i.UserID,
			&i.ID_2,
			&i.Name,
			&i.PollID_2,
			&i.Count,
			&i.CreatedAt_2,
			&i.UpdatedAt,
			&i.ID_3,
			&i.UserID_2,
			&i.Title,
			&i.Description,
			&i.Category,
			&i.CreatedAt_3,
			&i.UpdatedAt_2,
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
