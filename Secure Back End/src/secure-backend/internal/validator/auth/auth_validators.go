package validator

import (
	"database/sql"
	"errors"

	authdb "github.com/dvg1130/Portfolio/secure-backend/repo/auth_db"
)

// length check
func AuthLenCheck(username string, password string) error {
	if len(username) < 8 || len(password) < 8 {
		return errors.New("username and password must be at least 8 characters")
	}
	return nil
}

// existing user
func ExistingUser(db *sql.DB, username string) (bool, error) {
	var exists bool
	err := db.QueryRow(authdb.UserExists, username).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil

}
