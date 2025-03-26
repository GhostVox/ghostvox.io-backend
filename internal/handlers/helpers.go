package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GhostVox/ghostvox.io-backend/internal/database"
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

func getWinner(options []database.Option) string {
	currentWinner := ""
	currentCount := int32(0)
	for _, option := range options {
		if option.Count > currentCount {
			currentWinner = option.Name
			currentCount = option.Count
		}
		if option.Count == currentCount {
			currentWinner = "draw"
		}
	}
	return currentWinner
}
