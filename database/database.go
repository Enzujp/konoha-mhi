package database

import (
	"database/sql"
	"errors"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"os"
)


var DB *sql.DB

func ConnectDB() (*sql.DB, error){
	var err error
	godotenv.Load()
	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		return nil, errors.New("no database string provided")
	}
	DB, err = sql.Open("postgres", DB_URL)
	if err != nil {
		return nil, err
	}
	
	return DB, nil
}

func Init() error {
	_, err := ConnectDB()
	if err != nil {
		return err
	}
	return nil
}

func CloseDatabase() error {
	if DB != nil{
		return DB.Close()
	}
	return nil
}