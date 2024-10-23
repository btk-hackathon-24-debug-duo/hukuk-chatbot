package utils

import (
	"os"
	"time"

	"github.com/btk-hackathon-24-debug-duo/project-setup/internal/models"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	User models.User `json:"user"`
	jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("JWT_PRIV_KEY"))

func CreateJWTToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
