package auth

import (
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func TestParseAuthHeader(t *testing.T) {
	// Test cases
	tests := []struct {
		header string
		want   string
	}{
		{"Bearer token", "token"},
		{"Token token", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.header, func(t *testing.T) {
			got, _ := ParseAuthHeader(tt.header)

			if got != tt.want {
				t.Errorf("ParseAuthHeader(%q) = %q, want %q", tt.header, got, tt.want)
			}
		})
	}
}

// Test for a valid JWT.
func TestValidateJWT_Valid(t *testing.T) {
	jwtSecret := "test_secret"
	userID := uuid.New()
	role := "admin"
	picture := "https://example.com/avatar.jpg"
	// Generate a token with a 15-minute expiration.
	tokenStr, err := GenerateJWTAccessToken(userID, role, picture, jwtSecret, 15*time.Minute)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	claims, err := ValidateJWT(tokenStr, jwtSecret)
	if err != nil {
		t.Fatalf("expected token to be valid, got error: %v", err)
	}

	// Verify the claims.
	if claims.Subject != userID.String() {
		t.Errorf("expected userID %s, got %s", userID.String(), claims.Subject)
	}
	if claims.Role != role {
		t.Errorf("expected role %s, got %s", role, claims.Role)
	}
	if claims.Issuer != "GhostVox" {
		t.Errorf("expected issuer 'GhostVox', got %s", claims.Issuer)
	}
}

// Test for an expired token.
func TestValidateJWT_Expired(t *testing.T) {
	jwtSecret := "test_secret"
	userID := uuid.New()
	role := "user"
	picture := "https://example.com/avatar.jpg"
	// Generate a token that expired 1 minute ago.
	tokenStr, err := GenerateJWTAccessToken(userID, role, picture, jwtSecret, -1*time.Minute)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	_, err = ValidateJWT(tokenStr, jwtSecret)
	if err == nil {
		t.Fatalf("expected error for expired token, got nil")
	}

	expectedErr := "token expired"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

// Test validation with an incorrect secret key.
func TestValidateJWT_InvalidSecret(t *testing.T) {
	jwtSecret := "correct_secret"
	wrongSecret := "wrong_secret"
	userID := uuid.New()
	role := "user"
	picture := "https://example.com/avatar.jpg"

	tokenStr, err := GenerateJWTAccessToken(userID, role, picture, jwtSecret, 15*time.Minute)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	_, err = ValidateJWT(tokenStr, wrongSecret)
	if err == nil {
		t.Fatalf("expected error for token with invalid secret, got nil")
	}
	// The error message may not be exactly "invalid token", so we simply check for non-nil.
}

// Test for a token with an unexpected signing method.
func TestValidateJWT_UnexpectedSigningMethod(t *testing.T) {
	// Create a token with the SigningMethodNone.
	claims := CustomClaims{
		UserId: "user-wrong-method",
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "GhostVox",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	// Note: SigningMethodNone is unsafe and normally disabled.
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	// Allow "none" signature for testing purposes.
	tokenStr, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("unexpected error signing token with none method: %v", err)
	}

	_, err = ValidateJWT(tokenStr, "any_secret")
	if err == nil {
		t.Fatalf("expected error for token with unexpected signing method, got nil")
	}
	expectedErr := "failed to parse token: unexpected signing method"
	if !errors.Is(err, errors.New(expectedErr)) && err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
func TestGenerateJWTAccessToken(t *testing.T) {
	jwtSecret := "my_super_secret"
	userID := uuid.New()
	role := "admin"
	accessTokenDuration := 15 * time.Minute
	picture := "https://example.com/avatar.jpg"

	tokenStr, err := GenerateJWTAccessToken(userID, role, picture, jwtSecret, accessTokenDuration)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Parse the token back to verify claims
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		t.Fatalf("error parsing token: %v", err)
	}
	if !token.Valid {
		t.Fatalf("expected token to be valid")
	}
	UserIDString := userID.String()
	if err != nil {
		t.Fatalf("error parsing userID: %v", err)
	}
	// Verify custom claims
	if claims.Subject != UserIDString {
		t.Errorf("expected userID %s, got %s", userID, claims.Subject)
	}
	if claims.Role != role {
		t.Errorf("expected role %s, got %s", role, claims.Role)
	}
	if claims.Issuer != "GhostVox" {
		t.Errorf("expected issuer 'GhostVox', got %s", claims.Issuer)
	}

	// Check that ExpiresAt is roughly as expected.
	// We allow a one-minute margin on either side to account for processing delays.
	now := time.Now()
	expectedExpire := now.Add(accessTokenDuration)
	if claims.ExpiresAt.Time.Before(expectedExpire.Add(-1*time.Minute)) || claims.ExpiresAt.Time.After(expectedExpire.Add(1*time.Minute)) {
		t.Errorf("expected expiration around %v, got %v", expectedExpire, claims.ExpiresAt.Time)
	}
}

func TestGenerateJWTRefreshToken(t *testing.T) {
	tokenStr, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if tokenStr == "" {
		t.Fatalf("expected non-empty token string")
	}
	expectedLen := hex.EncodedLen(32) // 32 bytes * 2 = 64 hex characters

	if len(tokenStr) != expectedLen {
		t.Fatalf("expected token string length %d, got %d", expectedLen, len(tokenStr))
	}
}
