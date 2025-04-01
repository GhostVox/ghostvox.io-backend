package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
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

	userRecord, err := qtx.UpdateUserProfile(ctx, database.UpdateUserProfileParams{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  NullStringHelper(user.LastName),
		UserName:  NullStringHelper(user.UserName),
	})
	if err != nil {
		return "", database.User{}, errors.New("Failed to update user")
	}

	refreshTokenString, err := auth.GenerateRefreshToken()
	if err != nil {
		fmt.Println("error creating refresh token:", err)
		return "", database.User{}, errors.New("Failed add refresh token to db")
	}

	_, err = qtx.UpdateRefreshToken(ctx, database.UpdateRefreshTokenParams{
		UserID: userRecord.ID,
		Token:  refreshTokenString,
	})
	if err != nil {
		fmt.Println("error creating refresh token:", err)
		return "", database.User{}, errors.New("Failed add refresh token to db")
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("error committing transaction:", err)
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

func CreatePollWithOptions(ctx context.Context, cfg *config.APIConfig, poll poll, userUUID uuid.UUID) (err error) {
	tx, err := cfg.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := cfg.Queries.WithTx(tx)
	exp, err := strconv.Atoi(poll.ExpiresAt)
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(time.Duration(exp) * 24 * time.Hour) // write a reusable helper for this and test.

	pollRecord, err := qtx.CreatePoll(ctx, database.CreatePollParams{
		UserID:      userUUID,
		Title:       poll.Title,
		Description: poll.Description,
		Category:    poll.Category,
		ExpiresAt:   expiresAt,
		Status:      database.PollStatus("Active"),
	})
	if err != nil {
		return err
	}

	names := make([]string, len(poll.Options))
	for i, option := range poll.Options {
		names[i] = option.Name
	}
	_, err = qtx.CreateOptions(ctx, database.CreateOptionsParams{
		PollID:  pollRecord.ID,
		Column2: names,
	})
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func CreateVoteAndUpdateOptionCount(ctx context.Context, cfg *config.APIConfig, userID, optionID, pollID uuid.UUID) (vote database.Vote, err error) {
	tx, err := cfg.DB.Begin()
	if err != nil {
		return database.Vote{}, err
	}
	defer tx.Rollback()
	qtx := cfg.Queries.WithTx(tx)

	vote, err = qtx.CreateVote(ctx, database.CreateVoteParams{
		UserID:   userID,
		PollID:   pollID,
		OptionID: optionID,
	})
	if err != nil {
		return database.Vote{}, err
	}

	_, err = qtx.UpdateOptionCount(ctx, optionID)
	if err != nil {
		return database.Vote{}, err
	}

	err = tx.Commit()
	if err != nil {
		return database.Vote{}, err
	}

	return vote, nil
}
