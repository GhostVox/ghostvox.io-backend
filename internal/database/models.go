// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PollStatus string

const (
	PollStatusActive   PollStatus = "Active"
	PollStatusInactive PollStatus = "Inactive"
	PollStatusArchived PollStatus = "Archived"
)

func (e *PollStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PollStatus(s)
	case string:
		*e = PollStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PollStatus: %T", src)
	}
	return nil
}

type NullPollStatus struct {
	PollStatus PollStatus
	Valid      bool // Valid is true if PollStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPollStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PollStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PollStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPollStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PollStatus), nil
}

type Option struct {
	ID        uuid.UUID
	Name      string
	PollID    uuid.UUID
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Poll struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpiresAt   time.Time
	Status      PollStatus
}

type RefreshToken struct {
	Token     string
	UserID    uuid.UUID
	CreatedAt time.Time
	ExpiresAt time.Time
}

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Email          string
	FirstName      string
	LastName       sql.NullString
	HashedPassword sql.NullString
	Provider       sql.NullString
	ProviderID     sql.NullString
	Role           string
}

type Vote struct {
	ID        uuid.UUID
	PollID    uuid.UUID
	OptionID  uuid.UUID
	CreatedAt time.Time
	UserID    uuid.UUID
}
