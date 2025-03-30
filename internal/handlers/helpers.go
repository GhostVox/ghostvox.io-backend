package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/google/uuid"
)

func NullStringHelper(value interface{}) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}

	return sql.NullString{String: value.(string), Valid: true}
}

func getLimitAndOffset(r *http.Request) (limit, offset int, err error) {
	limitParam := r.URL.Query().Get("limit")
	if limitParam == "" {
		limitParam = "20"
	}

	offsetParam := r.URL.Query().Get("offset")
	if offsetParam == "" {
		offsetParam = "0"
	}

	limit, err = strconv.Atoi(limitParam)
	if err != nil {
		err = fmt.Errorf("Invalid limit parameter: %w", err)
		return 0, 0, err
	}

	offset, err = strconv.Atoi(offsetParam)
	if err != nil {
		err = fmt.Errorf("Invalid offset parameter: %w", err)
		return 0, 0, err
	}

	return limit, offset, nil
}

func getWinner(options []database.Option) uuid.UUID {
	currentWinner := uuid.Nil
	currentCount := int32(0)
	for _, option := range options {
		if option.Count > currentCount {
			currentWinner = option.ID
			currentCount = option.Count
			fmt.Println(option.Count)
			continue
		}
		if option.Count == currentCount {
			currentWinner = uuid.Nil
		}
	}
	return currentWinner
}
