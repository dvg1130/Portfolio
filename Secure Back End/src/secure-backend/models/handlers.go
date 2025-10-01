package models

import (
	"database/sql"
	"net/http"
)

type AdminHandlers struct {
	//admin
	AdminGetAll http.HandlerFunc
	AdminGetOne http.HandlerFunc
	AdminUpdate http.HandlerFunc
}

type AuthHandlers struct {
	//auth
	Handler  http.HandlerFunc
	Login    http.HandlerFunc
	Register http.HandlerFunc
	Logout   http.HandlerFunc
}

type DataHandlers struct {
	S *sql.DB

	//snake
	SnakeGetAll http.HandlerFunc
	SnakeGetOne http.HandlerFunc
	SnakePost   http.HandlerFunc
	SnakeUpdate http.HandlerFunc
	SnakeDelete http.HandlerFunc

	//feed
	SnakeFeedGet    http.HandlerFunc
	SnakeFeedPost   http.HandlerFunc
	SnakeFeedUpdate http.HandlerFunc
	SnakeFeedDelete http.HandlerFunc

	//health
	SnakeHealthGet    http.HandlerFunc
	SnakeHealthPost   http.HandlerFunc
	SnakeHealthUpdate http.HandlerFunc
	SnakeHealthDelete http.HandlerFunc
}

type BreedingHandlers struct {
	//breed
	SnakeBreedGetAll     http.HandlerFunc
	SnakeBreedGetBySnake http.HandlerFunc
	SnakeBreedGetOne     http.HandlerFunc
	SnakeBreedPost       http.HandlerFunc
	SnakeBreedUpdate     http.HandlerFunc
	SnakeBreedDelete     http.HandlerFunc
}
