package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/google/uuid"
)

func TestAdminRoleMiddleware(t *testing.T) {
	// Setup config
	apiConfig := config.APIConfig{
		GhostvoxSecretKey: "testsecretkey",
	}

	// Setup test handler that the middleware will wrap
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Create the middleware-wrapped handler
	adminHandler := AdminRole(&apiConfig, testHandler)

	t.Run("Valid Admin Token", func(t *testing.T) {
		// Generate admin token
		userID := uuid.New()
		role := "admin"
		picture := "https://example.com/avatar.jpg"
		firstName := "John"
		lastName := "Doe"
		email := "john.doe@example.com"
		duration := time.Duration(time.Hour * 24)

		token, err := auth.GenerateJWTAccessToken(
			userID, role, picture, firstName, lastName, email,
			apiConfig.GhostvoxSecretKey, duration)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Create test request with token
		req := httptest.NewRequest("GET", "/admin-resource", nil)
		req.AddCookie(&http.Cookie{
			Name:     "accessToken",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})

		// Record the response
		rec := httptest.NewRecorder()

		// Serve the request
		adminHandler.ServeHTTP(rec, req)

		// Check response
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}
	})

	t.Run("Non-Admin User Token", func(t *testing.T) {
		// Generate non-admin token
		userID := uuid.New()
		role := "user" // Non-admin role
		picture := "https://example.com/avatar.jpg"
		firstName := "Regular"
		lastName := "User"
		email := "regular.user@example.com"
		duration := time.Duration(3600 * time.Second)

		token, err := auth.GenerateJWTAccessToken(
			userID, role, picture, firstName, lastName, email, apiConfig.GhostvoxSecretKey, duration)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Create test request with token
		req := httptest.NewRequest("GET", "/admin-resource", nil)
		req.AddCookie(&http.Cookie{
			Name:     "accessToken",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})

		// Record the response
		rec := httptest.NewRecorder()

		// Serve the request
		adminHandler.ServeHTTP(rec, req)

		// Check response - should be forbidden for non-admin
		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status Forbidden, got %d", rec.Code)
		}
	})

	t.Run("No Token", func(t *testing.T) {
		// Create test request without token
		req := httptest.NewRequest("GET", "/admin-resource", nil)

		// Record the response
		rec := httptest.NewRecorder()

		// Serve the request
		adminHandler.ServeHTTP(rec, req)

		// Check response - should be unauthorized
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status Unauthorized, got %d", rec.Code)
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		// Create test request with invalid token
		req := httptest.NewRequest("GET", "/admin-resource", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken123")

		// Record the response
		rec := httptest.NewRecorder()

		// Serve the request
		adminHandler.ServeHTTP(rec, req)

		// Check response - should be unauthorized
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status Unauthorized, got %d", rec.Code)
		}
	})
}
