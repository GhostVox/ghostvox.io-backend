package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, password string) error {
	fmt.Printf("password in verify password: %s , hashed password in verify password: %s", password, hashedPassword)
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
