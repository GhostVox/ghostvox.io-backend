package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "supersecret123"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Ensure hashed password is not empty
	if hashedPassword == "" {
		t.Fatalf("expected hashed password, got empty string")
	}

	// Ensure hashed password is not the same as the original password
	if hashedPassword == password {
		t.Fatalf("hashed password should not match original password")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "supersecret123"

	// Hash the password first
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	// Test correct password verification
	err = VerifyPassword(hashedPassword, password)
	if err != nil {
		t.Fatalf("expected password verification to succeed, got error: %v", err)
	}

	// Test incorrect password verification
	wrongPassword := "wrongpassword"
	err = VerifyPassword(hashedPassword, wrongPassword)
	if err == nil {
		t.Fatalf("expected password verification to fail with wrong password, but it passed")
	}
}
