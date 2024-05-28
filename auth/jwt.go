package auth

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID      int    `json:"user_id"`
	AccessToken string `json:"access_token"`
	jwt.RegisteredClaims
}

func Issue(claims UserClaims) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("[WARN] JWT_SECRET not set, using default")
		secret = "default"
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(secret))
}

func Verify(token string) (*UserClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("[WARN] JWT_SECRET not set, using default")
		secret = "default"
	}

	var claims UserClaims

	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("JWT invalid")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("JWT expired")
	}

	return &claims, nil
}
