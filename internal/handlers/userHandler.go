package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
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
}

type UserHandler struct {
	cfg *config.APIConfig
}

func NewUserHandler(cfg *config.APIConfig) *UserHandler {
	return &UserHandler{
		cfg: cfg,
	}
}

// func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
// 	users, err := h.cfg.Queries.GetUsers(r.Context())
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			chooseError(w, http.StatusNotFound, err)
// 			return
// 		}
// 		chooseError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	respondWithJSON(w, http.StatusOK, users)
// 	return
// // }

// func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
// 	id := r.PathValue("userId")
// 	if id == "" {
// 		chooseError(w, http.StatusBadRequest, errors.New("missing id"))
// 		return
// 	}
// 	UserUUID, err := uuid.Parse(id)
// 	if err != nil {
// 		chooseError(w, http.StatusBadRequest, errors.New("invalid id"))
// 		return
// 	}
// 	user, err := h.cfg.Queries.GetUserById(r.Context(), UserUUID)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			chooseError(w, http.StatusNotFound, err)
// 			return
// 		}
// 		chooseError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	respondWithJSON(w, http.StatusOK, user)
// 	return
// }

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie, err := r.Cookie("access_token")
	if err != nil {
		chooseError(w, http.StatusUnauthorized, errors.New("missing access token"))
		return
	}
	claims, err := auth.ValidateJWT(accessTokenCookie.Value, h.cfg.GhostvoxSecretKey)
	if err != nil {
		chooseError(w, http.StatusUnauthorized, err)
		return
	}

	var user User
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		chooseError(w, http.StatusInternalServerError, err)
		return
	}
	UserUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		chooseError(w, http.StatusBadRequest, errors.New("invalid id"))
		return
	}
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		chooseError(w, http.StatusInternalServerError, fmt.Errorf("password hashing failed: %w", err))
		return
	}
	user.Password = hashedPassword
	user.ID = UserUUID
	refreshToken, updatedUserRecord, err := updateUserAndRefreshToken(r.Context(), h.cfg.DB, h.cfg.Queries, user)
	if err != nil {
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	accessToken, err := auth.GenerateJWTAccessToken(updatedUserRecord.ID, updatedUserRecord.Role, updatedUserRecord.PictureUrl.String, h.cfg.GhostvoxSecretKey, h.cfg.AccessTokenExp)
	if err != nil {
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	SetCookiesHelper(w, refreshToken, accessToken, h.cfg)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie, err := r.Cookie("access_token")
	if err != nil {
		chooseError(w, http.StatusUnauthorized, errors.New("missing access token"))
		return
	}
	claims, err := auth.ValidateJWT(accessTokenCookie.Value, h.cfg.GhostvoxSecretKey)
	if err != nil {
		chooseError(w, http.StatusUnauthorized, err)
		return
	}

	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		chooseError(w, http.StatusBadRequest, err)
		return
	}

	err = h.cfg.Queries.DeleteUser(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
	return
}
