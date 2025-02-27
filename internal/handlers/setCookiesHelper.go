package handlers

import (
	"net/http"
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/config"
)

func SetCookiesHelper(w http.ResponseWriter, refreshToken, accessToken string, cfg *config.APIConfig) {
	// Set cookies for the user's session
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   cfg.Mode == "Dev",
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Expires:  time.Now().Add(cfg.RefreshTokenExp),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   cfg.Mode == "Dev",
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Expires:  time.Now().Add(cfg.AccessTokenExp),
	})
	w.Header().Set("Authorization", "Bearer "+accessToken)

	// Respond with a success message
	if refreshToken == "" && accessToken == "" {
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "User logged out successfully"})
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})

}
