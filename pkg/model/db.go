package model

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	util "../util"
	"golang.org/x/crypto/bcrypt"
)

func (u *User) VerifyUserCredentials() (string, error) {
	db, err := util.CreateDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	var email string
	var hash string
	var verified bool

	query := `SELECT email, password, verified FROM users WHERE email = '` + u.Email + `';`
	row := db.DBClient.QueryRow(query)
	switch err := row.Scan(&email, &hash, &verified); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return "User does not exist", err
	case nil:
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Password))
		if err != nil {
			log.Println(err)
			return "", err
		}
		if verified {
			return "User is verified", nil
		}
		return "User is not verified", nil
	default:
		log.Println(err)
		return "", err
	}
}

func (u *User) AddUserToUserTable() (bool, error) {
	db, err := util.CreateDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	var email string

	query := `SELECT email FROM users WHERE email = '` + u.Email + `';`
	row := db.DBClient.QueryRow(query)
	switch err := row.Scan(&email); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
		if err != nil {
			log.Println(err)
			return false, err
		}

		query := fmt.Sprintf("INSERT INTO users (email, password, fullname, address, phonenumber, verified) Values ('%s', '%s', '%s', '%s', '%s', false)", u.Email, encryptedPassword, u.FullName, "", u.PhoneNumber)
		_, err = db.DBClient.Exec(query)
		if err != nil {
			log.Println(err)
			return false, err
		}
		return true, nil
	case nil:
		return false, errors.New("User already exists in table")
	default:
		log.Println(err)
		return false, err
	}
}

func (u *User) UpdateUserTableVerified() error {
	db, err := util.CreateDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	query := fmt.Sprintf("UPDATE users set verified = true WHERE email = '%s'", u.Email)
	_, err = db.DBClient.Exec(query)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (u *User) GetUserDetails() (User, error) {
	var user User
	db, err := util.CreateDBConnection()
	if err != nil {
		log.Println(err)
		return user, err
	}

	query := fmt.Sprintf("SELECT email, fullname, address, phonenumber FROM users WHERE email = '%s';", u.Email)
	log.Println(query)
	row := db.DBClient.QueryRow(query)
	switch err := row.Scan(&user.Email, &user.FullName, &user.Address, &user.PhoneNumber); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, err
	case nil:
		return user, nil
	default:
		log.Println(err)
		return user, err
	}
}
