package handlers

import (
	"database/sql"
)

func NullStringHelper(value interface{}) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}

	return sql.NullString{String: value.(string), Valid: true}
}
