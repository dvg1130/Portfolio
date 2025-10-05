package api

import (
	"net/http"

	"github.com/dvg1130/Portfolio/secure-auth/internal/middleware"
	validator "github.com/dvg1130/Portfolio/secure-auth/internal/validator/util"
	"github.com/dvg1130/Portfolio/secure-auth/models"
)

func InitRoutes_Auth(router *http.ServeMux, h *models.AuthHandlers) {

	//routes
	router.HandleFunc("/", h.Handler)
	router.Handle("/auth/login",
		middleware.RequestIDMiddleware(&h.Logger)(
			validator.Method(http.MethodPost,
				http.HandlerFunc(h.Login),
			)))
	router.HandleFunc("/auth/register", validator.Method(http.MethodPost, h.Register))
	router.HandleFunc("auth/logout", validator.Method(http.MethodPost, h.Logout))
	router.HandleFunc("/auth/token/refresh", validator.Method(http.MethodPost, h.TokenRefresh))
}
