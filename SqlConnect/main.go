package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// server struct
type Server struct {
	Router *http.ServeMux
	DB     *sql.DB
}

// create server and routes
func NewServer(db *sql.DB) *Server {
	s := &Server{
		Router: http.NewServeMux(),
		DB:     db,
	}
	// s.Router.Use(middleware.Logger)
	s.Router.HandleFunc("/", handler)
	s.Router.HandleFunc("POST /user", createUser)

	return s

}

// hander func
func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(" successfull connection to server!"))
}
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(" Successfull post to User Endpoint!"))
}

func main() {

	//init db
	db, err := DBClient()
	if err != nil {
		log.Fatal("failed to connect to db", err)
	}

	//init server
	server := NewServer(db)
	err = http.ListenAndServe(":8001", server.Router)
	if err != nil {
		fmt.Println("Error connecting to server")
	}

}
