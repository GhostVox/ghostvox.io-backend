package auth

import (
	"crypto/rand"
	"encoding/base64"
)

// generateRandomState creates a secure random string to use as an OAuth state parameter.
func GenerateRandomState() string {
	// Create a 32-byte random sequence
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		// Handle error safely (this should rarely happen)
		return "fallbackstate" // A less secure fallback (not recommended for production)
	}

	// Encode to base64 to make it URL-safe
	return base64.URLEncoding.EncodeToString(b)
}
