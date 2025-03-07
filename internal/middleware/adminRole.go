package middleware

import (
	"net/http"
	"strings"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
)

const accessTokenCookieName string = "accessToken"

func AdminRole(cfg *config.APIConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(accessTokenCookieName)
		if err != nil {
			http.Error(w, "access token missing", http.StatusUnauthorized)
			return
		}
		claims, err := auth.ValidateJWT(cookie.Value, cfg.GhostvoxSecretKey)
		if err != nil {
			http.Error(w, "invalid access token", http.StatusUnauthorized)
			return
		}
		if !strings.EqualFold(claims.Role, "admin") {
			http.Error(w, "unauthorized", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
