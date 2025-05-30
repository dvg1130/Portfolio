package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// server struct
type Server struct {
	Router *chi.Mux
}

// create server and routes
func NewServer() *Server {
	s := &Server{
		Router: chi.NewRouter(),
	}
	//logger
	s.Router.Use(middleware.Logger)
	s.Router.Get("/v1", handler2)

	return s
}

// hander func
func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Connection Successful"))
}

func handler2(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Connection Successful to v1"))
}

func main() {
	//init server
	server := NewServer()

	err := http.ListenAndServe(":8001", server.Router)
	if err != nil {
		fmt.Println("Error with connection", nil)
	}

}
