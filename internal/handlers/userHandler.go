package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
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
// 			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), err)
// 			return
// 		}
// 		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err)
// 		return
// 	}

// 	respondWithJSON(w, http.StatusOK, users)
// 	return
// // }

// func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
// 	id := r.PathValue("userId")
// 	if id == "" {
// 		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), errors.New("missing id"))
// 		return
// 	}
// 	UserUUID, err := uuid.Parse(id)
// 	if err != nil {
// 		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), errors.New("invalid id"))
// 		return
// 	}
// 	user, err := h.cfg.Queries.GetUserById(r.Context(), UserUUID)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), err)
// 			return
// 		}
// 		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal Server Error", err)
// 		return
// 	}

// 	respondWithJSON(w, http.StatusOK, user)
// 	return
// }

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie, err := r.Cookie("access_token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Missing Access Token", err)
		return
	}
	claims, err := auth.ValidateJWT(accessTokenCookie.Value, h.cfg.GhostvoxSecretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Unauthorized", err)
		return
	}

	var user User
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal Server Error", err)
		return
	}
	UserUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid ID", err)
		return
	}
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Password Hashing Failed", err)
		return
	}
	user.Password = hashedPassword
	user.ID = UserUUID
	refreshToken, updatedUserRecord, err := updateUserAndRefreshToken(r.Context(), h.cfg.DB, h.cfg.Queries, user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal Server Error", err)
		return
	}

	accessToken, err := auth.GenerateJWTAccessToken(updatedUserRecord.ID, updatedUserRecord.Role, updatedUserRecord.PictureUrl.String, h.cfg.GhostvoxSecretKey, h.cfg.AccessTokenExp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Token Generation Failed", err)
		return
	}

	SetCookiesHelper(w, refreshToken, accessToken, h.cfg)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie, err := r.Cookie("access_token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Missing access token", err)
		return
	}
	claims, err := auth.ValidateJWT(accessTokenCookie.Value, h.cfg.GhostvoxSecretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Unauthorized", err)
		return
	}

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

	respondWithJSON(w, http.StatusNoContent, nil)
	return
}
