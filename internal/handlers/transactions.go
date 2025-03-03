package handlers

import (
	"context"
	"database/sql"
	"errors"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/lib/pq"
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
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "23505" {
				return "", database.User{}, errors.New("Email already exists")
			}
		}
		return "", database.User{}, errors.New("Failed to create user")
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
