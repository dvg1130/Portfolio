package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//db variables

var (
	DATABASE_URL,
	DB_DRIVER,
	JWT_SECRET_KEY,
	PORT string
)

// .env variables
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Couldn't load env variables")
	}
	DATABASE_URL = os.Getenv("DATABASE_URL")
	DB_DRIVER = os.Getenv("DB_DRIVER")
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	PORT = os.Getenv("PORT")

}

//db client

func DBClient() (*sql.DB, error) {
	db, err := sql.Open(DB_DRIVER, DATABASE_URL)
	if err != nil {

		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Connected to sql db")
	return db, nil
}
