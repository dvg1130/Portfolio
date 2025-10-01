package api

import (
	"net/http"

	"github.com/dvg1130/Portfolio/secure-backend/internal/middleware"
	validator "github.com/dvg1130/Portfolio/secure-backend/internal/validator/util"
	"github.com/dvg1130/Portfolio/secure-backend/models"
)

func InitRoutes_Breeding(router *http.ServeMux, h *models.BreedingHandlers) {

	//routes
	//breeding routes
	router.Handle("/dashboard/breeding/all", middleware.AuthMiddleware(
		validator.Method(http.MethodGet,
			http.HandlerFunc(h.SnakeBreedGetAll)),
	),
	)
	router.Handle("/dashboard/breeding/snake", middleware.AuthMiddleware(
		validator.Method(http.MethodGet,
			http.HandlerFunc(h.SnakeBreedGetBySnake)),
	),
	)

	router.Handle("/dashboard/breeding/one", middleware.AuthMiddleware(
		validator.Method(http.MethodGet,
			http.HandlerFunc(h.SnakeBreedGetOne)),
	),
	)

	router.Handle("/dashboard/breeding/post", middleware.AuthMiddleware(
		validator.Method(http.MethodPost,
			http.HandlerFunc(h.SnakeBreedPost)),
	),
	)

	router.Handle("/dashboard/breeding/update", middleware.AuthMiddleware(
		validator.Method(http.MethodPatch,
			http.HandlerFunc(h.SnakeBreedUpdate)),
	),
	)

	router.Handle("/dashboard/breeding/delete", middleware.AuthMiddleware(
		validator.Method(http.MethodDelete,
			http.HandlerFunc(h.SnakeBreedDelete)),
	),
	)

}
