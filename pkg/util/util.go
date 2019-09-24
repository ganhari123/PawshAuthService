package util

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DBConnDetails struct {
	Host          string
	Port          int
	Username      string
	Database      string
	Password      string
	ConnectionURL string
}

type HttpResponse struct {
	Body    string
	Success bool
	Error   string
}

type Database struct {
	DBClient *sql.DB
}

const (
	HOST           = "localhost"
	PORT           = 5432
	USERNAME       = "hari"
	DATABASE       = "hari"
	PASSWORD       = ""
	CONNECTION_URL = "postgresql://localhost"
)

func CreateDBConnection() (Database, error) {
	var dbClient Database
	dbconn := DBConnDetails{
		Host:          HOST,
		Port:          PORT,
		Username:      USERNAME,
		Database:      DATABASE,
		Password:      PASSWORD,
		ConnectionURL: CONNECTION_URL,
	}

	var db *sql.DB

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", dbconn.Host, dbconn.Port, dbconn.Username, dbconn.Password, dbconn.Database)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return dbClient, err
	}

	fmt.Println("Connection made successful")
	dbClient.DBClient = db

	return dbClient, nil
}
