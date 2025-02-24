package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChooseError(t *testing.T) {
	tests := []struct {
		name         string
		code         int
		err          error
		expectedCode int
		expectedMsg  string
	}{
		{"Not Found", http.StatusNotFound, nil, http.StatusNotFound, "Resource not found"},
		{"Unauthorized", http.StatusUnauthorized, nil, http.StatusUnauthorized, "Resource access denied"},
		{"Bad Request", http.StatusBadRequest, nil, http.StatusBadRequest, "Bad request"},
		{"Method Not Allowed", http.StatusMethodNotAllowed, nil, http.StatusMethodNotAllowed, "Method not allowed"},
		{"Not Implemented", http.StatusNotImplemented, nil, http.StatusNotImplemented, "Not implemented"},
		{"Default Case", http.StatusInternalServerError, nil, http.StatusInternalServerError, "Internal Server Error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			chooseError(rec, tt.code, tt.err)

			resp := rec.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			var body map[string]string
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if body["error"] != tt.expectedMsg {
				t.Errorf("expected message %q, got %q", tt.expectedMsg, body["error"])
			}
		})
	}
}
