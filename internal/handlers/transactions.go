package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/google/uuid"
)

func addUserAndRefreshToken(ctx context.Context, db *sql.DB, queries *database.Queries, user *User) (string, database.User, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", database.User{}, err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)

	userRecord, err := qtx.CreateUser(ctx, database.CreateUserParams{
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       NullStringHelper(user.LastName),
		HashedPassword: NullStringHelper(user.Password),
		Provider:       NullStringHelper(user.Provider),
		ProviderID:     NullStringHelper(user.ProviderID),
		PictureUrl:     NullStringHelper(user.PictureURL),
		Role:           user.Role,
	})
	if err != nil {

		return "", database.User{}, fmt.Errorf("Failed to add user error: %v", err)
	}

	refreshTokenString, err := auth.GenerateRefreshToken()
	qtx.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		UserID: userRecord.ID,
		Token:  refreshTokenString,
	})
	if err != nil {

		return "", database.User{}, errors.New("Failed add refresh token to db")
	}

	err = tx.Commit()
	if err != nil {
		return "", database.User{}, errors.New("Failed to commit transaction")
	}

	return refreshTokenString, userRecord, nil
}

func updateUserAndRefreshToken(ctx context.Context, db *sql.DB, queries *database.Queries, user User) (string, database.User, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", database.User{}, err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)

	userRecord, err := qtx.UpdateUser(ctx, database.UpdateUserParams{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       NullStringHelper(user.LastName),
		HashedPassword: NullStringHelper(user.Password),
		Provider:       NullStringHelper(user.Provider),
		ProviderID:     NullStringHelper(user.ProviderID),
		PictureUrl:     NullStringHelper(user.PictureURL),
		Role:           user.Role,
	})
	if err != nil {
		return "", database.User{}, errors.New("Failed to update user")
	}

	refreshTokenString, err := auth.GenerateRefreshToken()
	qtx.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		UserID: userRecord.ID,
		Token:  refreshTokenString,
	})
	if err != nil {
		return "", database.User{}, errors.New("Failed add refresh token to db")
	}

	err = tx.Commit()
	if err != nil {
		return "", database.User{}, errors.New("Failed to commit transaction")
	}

	return refreshTokenString, userRecord, nil
}

func deleteAndReplaceRefreshToken(ctx context.Context, cfg *config.APIConfig, userID uuid.UUID) (string, error) {
	tx, err := cfg.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	qtx := cfg.Queries.WithTx(tx)

	err = qtx.DeleteRefreshTokenByUserID(ctx, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	refreshTokenString, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	refreshRecord, err := qtx.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		UserID:    userID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(cfg.RefreshTokenExp),
	})
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return refreshRecord.Token, nil
}

func CreatePollWithOptions(ctx context.Context, db *sql.DB, cfg *config.APIConfig, poll poll) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := cfg.Queries.WithTx(tx)

	expiresAt := time.Now().Add(time.Duration(poll.ExpiresAt))
	pollRecord, err := qtx.CreatePoll(ctx, database.CreatePollParams{
		UserID:      poll.UserID,
		Title:       poll.Title,
		Description: poll.Description,
		Category:    poll.Category,
		ExpiresAt:   expiresAt,
		Status:      database.PollStatus("Active"),
	})
	if err != nil {
		return err
	}
	for _, option := range poll.Options {
		_, err := qtx.CreateOption(ctx, database.CreateOptionParams{
			PollID: pollRecord.ID,
			Name:   option.Name,
		})
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
