package config

import (
	"log"
	"os"

	"github.com/dvg1130/Portfolio/secure-auth/models"
	"github.com/joho/godotenv"
)

//init & load values

var AuthConfig models.AuthConfigStruct

func init() {

	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal("error loading env file")
	}

	AuthConfig = models.AuthConfigStruct{
		DATABASE_URL:   os.Getenv("AUTH_DATABASE_URL"),
		DB_DRIVER:      os.Getenv("AUTH_DB_DRIVER"),
		PORT:           os.Getenv("AUTH_PORT"),
		JWT_SECRET_KEY: os.Getenv("JWT_SECRETt_KEY"),
		REDIS_ADDR:     os.Getenv("REDIS_ADDR"),
	}

}
