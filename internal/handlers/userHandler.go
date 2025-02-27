package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/google/uuid"
)

type User struct {
	Email        string `json:"email,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Password     string `json:"password,omitempty"`
	Provider     string `json:"provider,omitempty"`
	ProviderID   string `json:"provider_id,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Role         string `json:"role,omitempty"`
	PictureURL   string `json:"picture,omitempty"`
}

type UserHandler struct {
	cfg *config.APIConfig
}

func NewUserHandler(cfg *config.APIConfig) *UserHandler {
	return &UserHandler{
		cfg: cfg,
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.cfg.DB.GetUsers(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, users)
	return
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("userId")
	if id == "" {
		chooseError(w, http.StatusBadRequest, errors.New("missing id"))
		return
	}
	UserUUID, err := uuid.Parse(id)
	if err != nil {
		chooseError(w, http.StatusBadRequest, errors.New("invalid id"))
		return
	}
	user, err := h.cfg.DB.GetUserById(r.Context(), UserUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
	return
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("userId")
	if id == "" {
		chooseError(w, http.StatusBadRequest, errors.New("missing id"))
		return
	}

	var user User
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		chooseError(w, http.StatusInternalServerError, err)
		return
	}
	UserUUID, err := uuid.Parse(id)
	if err != nil {
		chooseError(w, http.StatusBadRequest, errors.New("invalid id"))
		return
	}
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		chooseError(w, http.StatusInternalServerError, fmt.Errorf("password hashing failed: %w", err))
		return
	}

	updatedUserRecord, err := h.cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             UserUUID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       NullStringHelper(user.LastName),
		HashedPassword: NullStringHelper(hashedPassword),
		Provider:       NullStringHelper(user.Provider),
		ProviderID:     NullStringHelper(user.ProviderID),
		PictureUrl:     NullStringHelper(user.PictureURL),
		Role:           user.Role,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusNotFound, err)
			return
		}
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		chooseError(w, http.StatusInternalServerError, err)
		return
	}
	h.cfg.DB.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		UserID: UserUUID,
		Token:  refreshToken,
	})

	accessToken, err := auth.GenerateJWTAccessToken(updatedUserRecord.ID, updatedUserRecord.Role, updatedUserRecord.PictureUrl.String, h.cfg.GhostvoxSecretKey, h.cfg.AccessTokenExp)
	if err != nil {
		chooseError(w, http.StatusInternalServerError, err)
		return
	}

	// Set Refresh Token as HTTP-only Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true, // Use HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(h.cfg.RefreshTokenExp), // 30 days
	})

	// Set Access Token as HTTP-only Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(h.cfg.AccessTokenExp), // 30 minutes
	})

	// Also send Access Token in the response header (optional)
	w.Header().Set("Authorization", "Bearer "+accessToken)

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
	return
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("userId")
	if id == "" {
		chooseError(w, http.StatusBadRequest, errors.New("missing id"))
		return
	}

	userUUID, err := uuid.Parse(id)
	if err != nil {
		chooseError(w, http.StatusBadRequest, err)
		return
	}

	err = h.cfg.DB.DeleteUser(r.Context(), userUUID)
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
