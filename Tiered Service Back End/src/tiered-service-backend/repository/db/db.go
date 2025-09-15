package db

import (
	"database/sql"
	"fmt"
	"tiered-service-backend/config"

	_ "github.com/go-sql-driver/mysql"
)

// db
func DBClient() (*sql.DB, error) {
	db, err := sql.Open(config.Config.DB_DRIVER, config.Config.DATABASE_URL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("connected to database")
	return db, nil
}
