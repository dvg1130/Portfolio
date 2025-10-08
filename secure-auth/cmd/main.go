package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dvg1130/Portfolio/secure-auth/internal/server"
	"github.com/dvg1130/Portfolio/secure-auth/logs"
	authdb "github.com/dvg1130/Portfolio/secure-auth/repo/auth_db"
)

// init server
func main() {

	auth_db, err := authdb.AuthDBClient()
	if err != nil {
		log.Fatal("failed to connect to auth db", err)

		defer auth_db.Close()
	}

	logger := logs.NewLogger()
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("failed to flush logger: %v\n", err)
		}
	}()

	server := server.AppServer(auth_db, logger)
	http.ListenAndServe(":8080", server.Router)
	if err != nil {
		fmt.Println("error starting server")
	}
}
