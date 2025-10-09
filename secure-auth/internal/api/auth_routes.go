package api

import (
	"net/http"

	"github.com/dvg1130/Portfolio/secure-auth/internal/middleware"
	validator "github.com/dvg1130/Portfolio/secure-auth/internal/validator/util"
	"github.com/dvg1130/Portfolio/secure-auth/models"
)

func InitRoutes_Auth(router *http.ServeMux, h *models.AuthHandlers) {

	router.Handle("/api/auth/login",
		middleware.RequestIDMiddleware(&h.Logger)(
			validator.Method(http.MethodPost,
				http.HandlerFunc(h.Login),
			)))
	router.HandleFunc("/auth/register", validator.Method(http.MethodPost, h.Register))
	router.HandleFunc("auth/logout", validator.Method(http.MethodPost, h.Logout))
	router.HandleFunc("/auth/token/refresh", validator.Method(http.MethodPost, h.TokenRefresh))
	router.HandleFunc("/health", h.Health)

}

func InitRoutes_Proxy(router *http.ServeMux, h *models.AuthHandlers) { // Internal only
	router.HandleFunc("/api/sdk/login", h.ProxyLogin)

}
