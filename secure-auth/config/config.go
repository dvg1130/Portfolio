package config

import (
	"log"
	"os"

	"github.com/dvg1130/Portfolio/secure-auth/models"
	"github.com/joho/godotenv"
)

//init & load values

var AuthConfig models.AuthConfigStruct
var ProxyConfig models.ProxyConfigStruct

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

	ProxyConfig = models.ProxyConfigStruct{
		BASE_URL: os.Getenv("BASE_URL"),

		AUTH_LOGIN_ROUTE: os.Getenv("AUTH_LOGIN_ROUTE"),
		SDK_LOGIN_ROUTE:  os.Getenv("SDK_LOGIN_ROUTE"),

		AUTH_REGISTER_ROUTE: os.Getenv("AUTH_REGISTER_ROUTE"),
		SDK_REGISTER_ROUTE:  os.Getenv("SDK_REGISTER_ROUTE"),

		AUTH_LOGOUT_ROUTE: os.Getenv("AUTH_LOGOUT_ROUTE"),
		SDK_LOGOUT_ROUTE:  os.Getenv("SDK_LOGOUT_ROUTE"),

		AUTH_REFRESH_TOKEN_ROUTE: os.Getenv("AUTH_REFRESH_TOKEN_ROUTE"),
		SDK_REFRESH_TOKEN_ROUTE:  os.Getenv("SDK_REFRESH_TOKEN_ROUTE"),
	}

}
