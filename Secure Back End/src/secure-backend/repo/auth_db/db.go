package authdb

import (
	"database/sql"
	"fmt"

	"github.com/dvg1130/Portfolio/secure-backend/config"
	_ "github.com/go-sql-driver/mysql"
)

func AuthDBClient() (*sql.DB, error) {
	db, err := sql.Open(config.AuthConfig.DB_DRIVER, config.AuthConfig.DATABASE_URL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("connected to auth db")
	return db, nil
}
