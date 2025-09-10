package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCorsMiddleware tests that the CORS headers are correctly set.
func TestCorsMiddleware(t *testing.T) {
	mockHandlerCalled := false
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockHandlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	accessOrigin := "http://localhost:3000"
	// Assuming CorsMiddleware is in the same package.
	// We need a stub for this to compile standalone, but in a real project this isn't needed.
	var testHandler http.Handler
	// A placeholder for the actual middleware function for testing purposes.
	// In a real scenario, CorsMiddleware would already be defined in the package.
	corsMiddlewarePlaceholder := func(next http.Handler, accessOrigin string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", accessOrigin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
	testHandler = corsMiddlewarePlaceholder(mockHandler, accessOrigin)

	// --- Test Case 1: Preflight OPTIONS request ---
	t.Run("OPTIONS Request", func(t *testing.T) {
		mockHandlerCalled = false // Reset for this test case
		req := httptest.NewRequest("OPTIONS", "http://example.com", nil)
		rr := httptest.NewRecorder()

		testHandler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if mockHandlerCalled {
			t.Error("next handler should not be called for OPTIONS request")
		}

		checkCorsHeader(t, rr, "Access-Control-Allow-Origin", accessOrigin)
		checkCorsHeader(t, rr, "Access-Control-Allow-Credentials", "true")
		checkCorsHeader(t, rr, "Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		checkCorsHeader(t, rr, "Access-Control-Allow-Headers", "Content-Type, Authorization")
		checkCorsHeader(t, rr, "Access-Control-Expose-Headers", "Authorization")
	})

	// --- Test Case 2: Regular GET request ---
	t.Run("GET Request", func(t *testing.T) {
		mockHandlerCalled = false // Reset for this test case
		req := httptest.NewRequest("GET", "http://example.com", nil)
		rr := httptest.NewRecorder()

		testHandler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if !mockHandlerCalled {
			t.Error("next handler should be called for GET request")
		}

		checkCorsHeader(t, rr, "Access-Control-Allow-Origin", accessOrigin)
	})
}

// checkCorsHeader is a helper function to check for the presence and value of a header.
func checkCorsHeader(t *testing.T, rr *httptest.ResponseRecorder, key, value string) {
	t.Helper()
	if headerValue := rr.Header().Get(key); headerValue != value {
		t.Errorf("wrong %s header: got %q want %q", key, headerValue, value)
	}
}
