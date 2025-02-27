package handlers

import (
	"context"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/google/uuid"
)

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
