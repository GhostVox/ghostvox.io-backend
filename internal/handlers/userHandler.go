package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Email      string    `json:"email,omitempty"`
	FirstName  string    `json:"first_name,omitempty"`
	LastName   string    `json:"last_name,omitempty"`
	Password   string    `json:"password,omitempty"`
	Provider   string    `json:"provider,omitempty"`
	ProviderID string    `json:"provider_id,omitempty"`
	Role       string    `json:"role,omitempty"`
	PictureURL string    `json:"picture,omitempty"`
	UserName   string    `json:"user_name,omitempty"`
}

type UserHandler struct {
	cfg       *config.APIConfig
	s3Handler *AWSS3Handler
}

func NewUserHandler(cfg *config.APIConfig, s3Handler *AWSS3Handler) *UserHandler {
	return &UserHandler{
		cfg:       cfg,
		s3Handler: s3Handler,
	}
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, claims *auth.CustomClaims) {

	var user User
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal Server Error", err)
		return
	}

	UserUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid ID", err)
		return
	}

	user.ID = UserUUID
	fmt.Println(user)
	refreshToken, updatedUserRecord, err := updateUserAndRefreshToken(r.Context(), h.cfg.DB, h.cfg.Queries, user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal Server Error", err)
		return
	}

	claimsData := auth.TokenClaimsData{
		UserID:    updatedUserRecord.ID,
		Role:      updatedUserRecord.Role,
		Picture:   updatedUserRecord.PictureUrl.String,
		FirstName: updatedUserRecord.FirstName,
		LastName:  updatedUserRecord.LastName.String,
		Email:     updatedUserRecord.Email,
		UserName:  updatedUserRecord.UserName.String,
	}
	accessToken, err := auth.GenerateJWTAccessToken(claimsData, h.cfg.GhostvoxSecretKey, h.cfg.AccessTokenExp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Token Generation Failed", err)
		return
	}

	SetCookiesHelper(w, http.StatusOK, refreshToken, accessToken, h.cfg)
	respondWithJSON(w, http.StatusOK, map[string]any{"msg": "User updated successfully"})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, claims *auth.CustomClaims) {

	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid user ID", err)
		return
	}

	err = h.cfg.Queries.DeleteUser(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "User not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal Server Error", err)
		return
	}
	fmt.Println("Deleteing user, ", userUUID)

	h.s3Handler.DeleteUserAvatar(w, r, claims)

}

func (h *UserHandler) GetUserStats(w http.ResponseWriter, r *http.Request, claims *auth.CustomClaims) {
	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid user ID", err)
		return
	}

	stats, err := h.cfg.Queries.GetUserStats(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "User not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal Server Error", err)
		return
	}

	respondWithJSON(w, http.StatusOK, stats)
}

func (h *UserHandler) AddUserName(w http.ResponseWriter, r *http.Request, claims *auth.CustomClaims) {
	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid user ID", err)
		return
	}

	var username string
	if err := json.NewDecoder(r.Body).Decode(&username); err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid request body", err)
		return
	}

	if len(username) < 3 {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Username must be at least 3 characters long", err)
		return
	}

	usernameRegex := regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	if !usernameRegex.MatchString(username) {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Username can only contain letters, numbers, underscores and hyphens", nil)
		return
	}

	exists, err := h.cfg.Queries.CheckUserNameExists(r.Context(), NullStringHelper(username))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal Server Error", err)
		return
	}

	if exists {
		respondWithError(w, http.StatusConflict, http.StatusText(http.StatusConflict), "Username already exists", err)
		return
	}
	user, err := h.cfg.Queries.UpdateUserName(r.Context(), database.UpdateUserNameParams{
		ID:       userUUID,
		UserName: NullStringHelper(username),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), "User not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal Server Error", err)
		return
	}

	claimsData := auth.TokenClaimsData{
		UserID:    user.ID,
		Role:      user.Role,
		Picture:   user.PictureUrl.String,
		FirstName: user.FirstName,
		LastName:  user.LastName.String,
		Email:     user.Email,
		UserName:  user.UserName.String,
	}

	token, err := auth.GenerateJWTAccessToken(claimsData, h.cfg.GhostvoxSecretKey, h.cfg.AccessTokenExp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal Server Error", err)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal Server Error", err)
		return
	}

	refreshRecord, err := h.cfg.Queries.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(h.cfg.RefreshTokenExp),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal Server Error", err)
		return
	}

	SetCookiesHelper(w, http.StatusOK, refreshRecord.Token, token, h.cfg)
	respondWithJSON(w, http.StatusOK, struct{ message string }{message: "User updated successfully"})
}
