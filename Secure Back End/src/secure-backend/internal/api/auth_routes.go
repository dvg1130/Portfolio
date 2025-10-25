package api

import (
	"net/http"

	"github.com/dvg1130/Portfolio/secure-backend/internal/middleware"
	validator "github.com/dvg1130/Portfolio/secure-backend/internal/validator/util"
	"github.com/dvg1130/Portfolio/secure-backend/models"
)

func InitRoutes_Auth(router *http.ServeMux, h *models.AuthHandlers) {

	//routes
	// router.HandleFunc("/", h.Handler)
	router.HandleFunc("/auth/login", middleware.WithCORS(validator.Method(http.MethodPost, h.LoginHandler)))
	router.HandleFunc("/auth/register", validator.Method(http.MethodPost, h.RegisterHandler))
	router.HandleFunc("/auth/logout", validator.Method(http.MethodPost, h.LogoutHandler))
	// router.HandleFunc("/token/refresh", validator.Method(http.MethodPost, h.TokenRefresh))

	router.Handle("/auth/token/refresh", middleware.AuthMiddleware(
		validator.Method(http.MethodPost,
			http.HandlerFunc(h.TokenRefresh)),
	),
	)

}
