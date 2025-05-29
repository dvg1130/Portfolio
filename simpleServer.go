package main

import (
	"fmt"
	"net/http"
)

// server struct
type Server struct {
	Router *http.ServeMux
}

// create server and routes
func NewServer() *Server {
	s := &Server{
		Router: http.NewServeMux(),
	}

	s.Router.HandleFunc("/", handler)

	return s
}

// hander func
func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Weecome to the Server"))

}

func main() {
	//init server
	server := NewServer()

	err := http.ListenAndServe(":8001", server.Router)
	if err != nil {
		fmt.Println("Issue with connection", err)
	}
}
