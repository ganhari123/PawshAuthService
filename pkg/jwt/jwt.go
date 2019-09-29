package jwt

import (
	"log"
	"time"

	model "../model"
	jwt "github.com/dgrijalva/jwt-go"
)

// SharedSecret - needs to be moved to a config file
var SharedSecret = []byte("")

// GenerateJwtAccessToken - Generates a jwt token and encodes with secret
func GenerateJwtAccessToken(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claim := model.JwtClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claim)
	tokenString, err := token.SignedString(SharedSecret)
	if err != nil {
		log.Println("ERROR")
		return "", err
	}
	return tokenString, nil
}
