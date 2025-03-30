package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type CustomClaims struct {
	Role       string `json:"role"`
	PictureUrl string `json:"picture_url"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	UserName   string `json:"user_name"`

	jwt.RegisteredClaims
}

func GenerateJWTAccessToken(userID uuid.UUID, role, pictureUrl, firstName, lastName, email, userName string, jwtSecretKey string, AccessTokenExpiresAt time.Duration) (string, error) {
	claims := CustomClaims{

		Role:       role,
		PictureUrl: pictureUrl,
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		UserName:   userName,

		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "GhostVox",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiresAt)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return tokenString, nil
}

func GenerateRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshToken := hex.EncodeToString(token)
	return refreshToken, nil

}
