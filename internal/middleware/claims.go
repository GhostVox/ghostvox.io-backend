package middleware

import (
	"context"
	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"net/http"
)

type contextKey string

// ClaimsMiddleWare is a middleware that validates the JWT from the accessToken cookie
// and adds the claims to the request context

const claimsKey contextKey = "user_claims"

func Authenticator(secretKey string) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("accessToken")
			if err != nil {
				http.Error(w, "access token missing", http.StatusUnauthorized)
				return
			}

			claims, err := auth.ValidateJWT(cookie.Value, secretKey)
			if err != nil {
				http.Error(w, "invalid access token", http.StatusUnauthorized)
				return

			}
			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func ClaimsFromContext(ctx context.Context) (*auth.CustomClaims, bool) {
	claims, ok := ctx.Value(claimsKey).(*auth.CustomClaims)
	return claims, ok
}
