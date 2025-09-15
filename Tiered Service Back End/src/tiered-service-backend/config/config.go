package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// config struct
type ConfigStruct struct {
	DATABASE_URL   string
	DB_DRIVER      string
	JWT_SECRET_KEY string
	PORT           string
	REDIS_ADDR     string
}

// init & load vars
var Config ConfigStruct

func init() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal("error loading env file")
	}

	Config = ConfigStruct{
		DATABASE_URL:   os.Getenv("DATABASE_URL"),
		DB_DRIVER:      os.Getenv("DB_DRIVER"),
		PORT:           os.Getenv("PORT"),
		JWT_SECRET_KEY: os.Getenv("JWT_SECRET_KEY"),
		REDIS_ADDR:     os.Getenv("REDIS_ADDR"),
	}

}
