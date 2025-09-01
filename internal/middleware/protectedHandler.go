package middleware

import (
	"net/http"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
)

type ProtectedHandler func(w http.ResponseWriter, r *http.Request, claims *auth.CustomClaims)

func (fn ProtectedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	fn(w, r, claims)
}
