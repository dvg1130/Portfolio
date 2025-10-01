package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dvg1130/Portfolio/secure-backend/internal/server"
	"github.com/dvg1130/Portfolio/secure-backend/logs"
	authdb "github.com/dvg1130/Portfolio/secure-backend/repo/auth_db"
	datadb "github.com/dvg1130/Portfolio/secure-backend/repo/data_db"
)

// init server
func main() {

	auth_db, err := authdb.AuthDBClient()
	if err != nil {
		log.Fatal("failed to connect to auth db", err)

		defer auth_db.Close()
	}

	data_db, err := datadb.DataDBClient()
	if err != nil {
		log.Fatal("failed to connect to data db", err)
	}

	logger := logs.NewLogger()
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("failed to flush logger: %v\n", err)
		}
	}()

	server := server.AppServer(auth_db, data_db, logger)
	http.ListenAndServe(":8003", server.Router)
	if err != nil {
		fmt.Println("error starting server")
	}
}
