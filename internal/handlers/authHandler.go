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
		chooseError(w, http.StatusBadRequest, fmt.Errorf("invalid request payload: %w", err))
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		chooseError(w, http.StatusInternalServerError, fmt.Errorf("password hashing failed: %w", err))
		return
	}

	// Generate refresh token
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		chooseError(w, http.StatusInternalServerError, fmt.Errorf("refresh token generation failed: %w", err))
		return
	}

	// Create user in the database
	userRecord, err := h.cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       NullStringHelper(user.LastName),
		HashedPassword: NullStringHelper(hashedPassword),
		Provider:       NullStringHelper(user.Provider),
		ProviderID:     NullStringHelper(user.ProviderID),
		Role:           user.Role,
	})
	if err != nil {
		chooseError(w, http.StatusInternalServerError, fmt.Errorf("database user creation failed: %w", err))
		return
	}
	_, err = h.cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    userRecord.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(h.cfg.RefreshTokenExp),
	})
	if err != nil {
		chooseError(w, http.StatusInternalServerError, fmt.Errorf("database refresh token creation failed: %w", err))
		return
	}

	// Generate Access Token (JWT)
	accessToken, err := auth.GenerateJWTAccessToken(userRecord.ID, userRecord.Role, h.cfg.GhostvoxSecretKey, h.cfg.AccessTokenExp)
	if err != nil {
		chooseError(w, http.StatusInternalServerError, fmt.Errorf("access token generation failed: %w", err))
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

	// Respond with a success message
	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login Login
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		chooseError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	userRecord, err := h.cfg.DB.GetUserByEmail(r.Context(), login.Email)
	if err != nil {
		chooseError(w, http.StatusUnauthorized, fmt.Errorf("invalid credentials: %w", err))
		return
	}

	if err := auth.VerifyPassword(login.Password, userRecord.HashedPassword.String); err != nil {
		chooseError(w, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}

	// Generate Refresh Token (JWT)
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		chooseError(w, http.StatusInternalServerError, fmt.Errorf("refresh token generation failed: %w", err))
		return
	}

	// Create Refresh Token in the database
	_, err = h.cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    userRecord.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(h.cfg.RefreshTokenExp),
	})
	if err != nil {
		chooseError(w, http.StatusInternalServerError, fmt.Errorf("database refresh token creation failed: %w", err))
		return
	}

	// Generate Access Token (JWT)
	accessToken, err := auth.GenerateJWTAccessToken(userRecord.ID, userRecord.Role, h.cfg.GhostvoxSecretKey, h.cfg.AccessTokenExp)
	if err != nil {
		chooseError(w, http.StatusInternalServerError, fmt.Errorf("access token generation failed: %w", err))
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

	// Respond with a success message
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "User logged in successfully"})

}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		chooseError(w, http.StatusUnauthorized, fmt.Errorf("refresh token not found"))
		return
	}

	refreshTokenRecord, err := h.cfg.DB.GetRefreshToken(r.Context(), refreshCookie.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusUnauthorized, fmt.Errorf("invalid refresh token"))
		}
		chooseError(w, http.StatusInternalServerError, errors.New("Failed to get refresh token record"))
		return
	}

	userRecord, err := h.cfg.DB.GetUserById(r.Context(), refreshTokenRecord.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			chooseError(w, http.StatusUnauthorized, fmt.Errorf("invalid user"))
		}
		chooseError(w, http.StatusInternalServerError, errors.New("Failed to get user record"))
		return
	}

	// Generate New Access Token
	accessToken, err := auth.GenerateJWTAccessToken(userRecord.ID, userRecord.Role, h.cfg.GhostvoxSecretKey, h.cfg.AccessTokenExp)
	if err != nil {
		chooseError(w, http.StatusInternalServerError, fmt.Errorf("failed to generate access token: %w", err))
		return
	}
	newRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		chooseError(w, http.StatusInternalServerError, fmt.Errorf("failed to generate refresh token: %w", err))
		return
	}

	newRefreshRecord, err := h.cfg.DB.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		UserID: userRecord.ID,
		Token:  newRefreshToken,
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

	// Set Refresh Token as HTTP-only Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshRecord.Token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(h.cfg.RefreshTokenExp), // 30 days
	})

	// Also send Access Token in the response header (optional)
	w.Header().Set("Authorization", "Bearer "+accessToken)

	// Respond with a success message
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Access token refreshed successfully"})

}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get Refresh Token from Cookie
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		chooseError(w, http.StatusUnauthorized, fmt.Errorf("refresh token not found"))
		return
	}

	// Delete Refresh Token from Database
	err = h.cfg.DB.DeleteRefreshToken(r.Context(), refreshToken.Value)
	if err != nil {
		chooseError(w, http.StatusInternalServerError, fmt.Errorf("database refresh token deletion failed: %w", err))
		return
	}

	// Clear Cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour), // Expire immediately
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour), // Expire immediately
	})

	// Respond with a success message
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "User logged out successfully"})
}
