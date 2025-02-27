// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO
    users (id,email, first_name, last_name, hashed_password,provider,provider_id,role)
VALUES
    ($1, $2, $3, $4, $5,$6,$7,$8)
RETURNING
    id, created_at, updated_at, email, first_name, last_name, hashed_password, provider, provider_id, role
`

type CreateUserParams struct {
	ID             uuid.UUID
	Email          string
	FirstName      string
	LastName       sql.NullString
	HashedPassword sql.NullString
	Provider       sql.NullString
	ProviderID     sql.NullString
	Role           string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.FirstName,
		arg.LastName,
		arg.HashedPassword,
		arg.Provider,
		arg.ProviderID,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.HashedPassword,
		&i.Provider,
		&i.ProviderID,
		&i.Role,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM
    users
WHERE
    id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT
    id, created_at, updated_at, email, first_name, last_name, hashed_password, provider, provider_id, role
FROM
    users
WHERE
    email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.HashedPassword,
		&i.Provider,
		&i.ProviderID,
		&i.Role,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT
    id, created_at, updated_at, email, first_name, last_name, hashed_password, provider, provider_id, role
FROM
    users
WHERE
    id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.HashedPassword,
		&i.Provider,
		&i.ProviderID,
		&i.Role,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
Select
    id, created_at, updated_at, email, first_name, last_name, hashed_password, provider, provider_id, role
FROM
    users
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Email,
			&i.FirstName,
			&i.LastName,
			&i.HashedPassword,
			&i.Provider,
			&i.ProviderID,
			&i.Role,
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

const updateUser = `-- name: UpdateUser :one
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
    updated_at = NOW()
WHERE id = $8 RETURNING id, created_at, updated_at, email, first_name, last_name, hashed_password, provider, provider_id, role
`

type UpdateUserParams struct {
	Email          string
	FirstName      string
	LastName       sql.NullString
	HashedPassword sql.NullString
	Provider       sql.NullString
	ProviderID     sql.NullString
	Role           string
	ID             uuid.UUID
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.Email,
		arg.FirstName,
		arg.LastName,
		arg.HashedPassword,
		arg.Provider,
		arg.ProviderID,
		arg.Role,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.HashedPassword,
		&i.Provider,
		&i.ProviderID,
		&i.Role,
	)
	return i, err
}
