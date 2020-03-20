package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/nikvas0/dc-homework/auth/objects"
)

const accessExpirationTime = 15 * time.Minute

func GenerateTokens(user *objects.User) (string, string, error) {
	token := objects.Token{}
	token.UserID = user.ID
	token.Email = user.Email
	token.ExpiresAt = time.Now().Add(accessExpirationTime).Unix()
	token.Issuer = "auth"

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
