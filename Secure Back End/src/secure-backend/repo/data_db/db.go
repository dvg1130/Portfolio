package datadb

import (
	"database/sql"
	"fmt"

	"github.com/dvg1130/Portfolio/secure-backend/config"
)

func DataDBClient() (*sql.DB, error) {
	db, err := sql.Open(config.DataConfig.DB_DRIVER, config.DataConfig.DATABASE_URL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("connected to data database")
	return db, nil

}
