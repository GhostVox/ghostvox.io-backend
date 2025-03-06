package handlers

import (
	"net/http"
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/config"
)

func SetCookiesHelper(w http.ResponseWriter, code int, refreshToken, accessToken string, cfg *config.APIConfig) {
	// Set cookies for the user's session
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Expires:  time.Now().Add(cfg.RefreshTokenExp),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Expires:  time.Now().Add(cfg.AccessTokenExp),
	})
	w.Header().Set("Authorization", "Bearer "+accessToken)

}
