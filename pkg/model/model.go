package model

import (
	"database/sql"
	"fmt"
	"log"

	util "../util"
	jwt "github.com/dgrijalva/jwt-go"
	bcrypt "golang.org/x/crypto/bcrypt"
)

type User struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

type JwtClaim struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func (u *User) VerifyUserCredentials() (bool, error) {
	db, err := util.CreateDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	var email string
	var hash string

	query := `SELECT email, password FROM users WHERE email = '` + u.Email + `';`
	row := db.DBClient.QueryRow(query)
	switch err := row.Scan(&email, &hash); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return false, err
	case nil:
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Password))
		if err != nil {
			log.Println(err)
			return false, err
		}
		return true, nil
	default:
		log.Println(err)
		return false, err
	}
}
