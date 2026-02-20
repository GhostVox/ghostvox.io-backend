package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
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

func getWinner(options []Option) string {
	currentWinner := ""
	currentCount := int32(0)
	for _, option := range options {
		if option.Count > currentCount {
			currentWinner = option.ID
			currentCount = option.Count
			continue
		}
		if option.Count == currentCount {
			currentWinner = ""
		}
	}
	return currentWinner
}

func SetCookiesHelper(w http.ResponseWriter, code int, refreshToken, accessToken string, cfg *config.APIConfig) {
	// Set cookies for the user's session
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Domain:   ".ghostvox.io",
		Expires:  time.Now().Add(cfg.RefreshTokenExp),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Domain:   ".ghostvox.io",
		Expires:  time.Now().Add(cfg.AccessTokenExp),
	})
	w.Header().Set("Authorization", "Bearer "+accessToken)

}

func AddRefreshToken(ctx context.Context, userID uuid.UUID, db *database.Queries) (string, error) {
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", err
	}
	_, err = db.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		UserID: userID,
		Token:  refreshToken,
	})
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}
