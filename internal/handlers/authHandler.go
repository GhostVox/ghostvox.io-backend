package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"
)

type AuthHandler struct {
	cfg *config.APIConfig
}
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(cfg *config.APIConfig) *AuthHandler {
	return &AuthHandler{
		cfg: cfg,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var user User
	defer r.Body.Close()

	// Decode JSON body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid request payload", err)
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Password hashing failed", err)
		return
	}
	user.Password = hashedPassword
	refreshToken, userRecord, err := addUserAndRefreshToken(r.Context(), h.cfg.DB, h.cfg.Queries, &user)
	if err != nil {
		if err.Error() == "Email already exists" {
			respondWithError(w, http.StatusConflict, "email", "Email already exists", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "User creation failed", err)
		return
	}

	// Generate Access Token (JWT)
	accessToken, err := auth.GenerateJWTAccessToken(userRecord.ID, userRecord.Role, userRecord.PictureUrl.String, h.cfg.GhostvoxSecretKey, h.cfg.AccessTokenExp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Access token generation failed", err)
		return
	}

	SetCookiesHelper(w, refreshToken, accessToken, h.cfg)

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login Login
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid request payload", err)
		return
	}

	userRecord, err := h.cfg.Queries.GetUserByEmail(r.Context(), login.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Invalid credentials", err)
		return
	}

	if err := auth.VerifyPassword(login.Password, userRecord.HashedPassword.String); err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Invalid credentials", err)
		return
	}

	// Generate Refresh Token (JWT)
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Refresh token generation failed", err)
		return
	}

	// Create Refresh Token in the database
	_, err = h.cfg.Queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    userRecord.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(h.cfg.RefreshTokenExp),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Database refresh token creation failed", err)
		return
	}

	// Generate Access Token (JWT)
	accessToken, err := auth.GenerateJWTAccessToken(userRecord.ID, userRecord.Role, userRecord.PictureUrl.String, h.cfg.GhostvoxSecretKey, h.cfg.AccessTokenExp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Access token generation failed", err)
		return
	}

	SetCookiesHelper(w, refreshToken, accessToken, h.cfg)

}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Refresh token not found", err)
		return
	}

	refreshTokenRecord, err := h.cfg.Queries.GetRefreshToken(r.Context(), refreshCookie.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Invalid refresh token", err)
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Failed to get refresh token record", err)
		return
	}

	userRecord, err := h.cfg.Queries.GetUserById(r.Context(), refreshTokenRecord.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Invalid user", err)
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Failed to get user record", err)
		return
	}

	// Generate New Access Token
	accessToken, err := auth.GenerateJWTAccessToken(userRecord.ID, userRecord.Role, userRecord.PictureUrl.String, h.cfg.GhostvoxSecretKey, h.cfg.AccessTokenExp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Failed to generate access token", err)
		return
	}
	newRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Failed to generate refresh token", err)
		return
	}

	newRefreshRecord, err := h.cfg.Queries.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		UserID: userRecord.ID,
		Token:  newRefreshToken,
	})

	SetCookiesHelper(w, newRefreshRecord.Token, accessToken, h.cfg)

}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get Refresh Token from Cookie
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid refresh token", err)
		return
	}

	// Delete Refresh Token from Database
	err = h.cfg.Queries.DeleteRefreshToken(r.Context(), refreshToken.Value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Failed to delete refresh token", err)
		return
	}

	// Clear Cookies
	SetCookiesHelper(w, "", "", h.cfg)

}
