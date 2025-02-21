package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Info("request started",
			"path", r.URL.Path,
			"method", r.Method)

		next(w, r)

		slog.Info("request completed",
			"path", r.URL.Path,
			"duration", time.Since(start))
	}
}
