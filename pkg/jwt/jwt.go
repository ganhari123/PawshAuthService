package jwt

import (
	"fmt"
	"log"
	"time"

	model "../model"
	jwt "github.com/dgrijalva/jwt-go"
)

var SharedSecret = []byte("secret_key")

func GenerateJwtToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claim := model.JwtClaim{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claim)
	tokenString, err := token.SignedString(SharedSecret)
	if err != nil {
		log.Println("ERROR")
		log.Fatal(err)
		return "", err
	}
	fmt.Println(tokenString)
	return tokenString, nil
}
