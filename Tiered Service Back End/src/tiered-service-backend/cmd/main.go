package main

import (
	"fmt"
	"log"
	"net/http"
	"tiered-service-backend/internal/server"
	"tiered-service-backend/repository"
	"tiered-service-backend/repository/db"
)

// init server
func main() {
	logger := repository.NewLogger()
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("failed to flush logger: %v\n", err)
		}
	}()

	database, err := db.DBClient()
	if err != nil {
		log.Fatal("failed to connectto database", err)

		defer database.Close()

	}
	server := server.NewServer(database, logger)
	http.ListenAndServe(":8001", server.Router)
	if err != nil {
		fmt.Println("Error starting server")
	}

}
