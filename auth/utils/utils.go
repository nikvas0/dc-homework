package utils

import (
	"os"
	"time"

	"auth/objects"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const accessExpirationTime = 15 * time.Minute
const confirmExpirationTime = 24 * time.Hour

func GenerateTokens(user *objects.User) (string, string, error) {
	token := objects.Token{}
	token.UserID = user.ID
	token.Email = user.Email
	token.Role = user.Role
	token.ExpiresAt = time.Now().Add(accessExpirationTime).Unix()
	token.Issuer = "auth-access"

	accessToken, err := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), token).SignedString([]byte(os.Getenv("ACCESS_TOKEN_KEY")))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uuid.NewRandom()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken.String(), nil
}

func GenerateConfirmToken() (string, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return token.String(), nil
}
