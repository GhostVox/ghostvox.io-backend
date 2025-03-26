// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: comments.sql

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const getTotalComments = `-- name: GetTotalComments :one
SELECT COUNT(*) FROM comments WHERE poll_id = $1
`

func (q *Queries) GetTotalComments(ctx context.Context, pollID uuid.UUID) (int64, error) {
	row := q.db.QueryRowContext(ctx, getTotalComments, pollID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getTotalCommentsByPollIDs = `-- name: GetTotalCommentsByPollIDs :many
SELECT poll_id, COUNT(*) as count
FROM comments
WHERE poll_id = ANY($1::uuid[])
GROUP BY poll_id
`

type GetTotalCommentsByPollIDsRow struct {
	PollID uuid.UUID
	Count  int64
}

func (q *Queries) GetTotalCommentsByPollIDs(ctx context.Context, dollar_1 []uuid.UUID) ([]GetTotalCommentsByPollIDsRow, error) {
	rows, err := q.db.QueryContext(ctx, getTotalCommentsByPollIDs, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTotalCommentsByPollIDsRow
	for rows.Next() {
		var i GetTotalCommentsByPollIDsRow
		if err := rows.Scan(&i.PollID, &i.Count); err != nil {
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
