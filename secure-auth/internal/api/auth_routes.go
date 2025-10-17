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

	router.HandleFunc("/api/auth/register", validator.Method(http.MethodPost, h.Register))

	router.Handle("/api/auth/logout",
		middleware.RequestIDMiddleware(&h.Logger)(
			middleware.AuthMiddleware(
				validator.Method(http.MethodPost,
					http.HandlerFunc(h.Logout)))))

	router.Handle("/api/auth/token/refresh",
		middleware.RequestIDMiddleware(&h.Logger)(
			middleware.AuthMiddleware(
				validator.Method(http.MethodPost,
					http.HandlerFunc(h.TokenRefresh)))))

	router.HandleFunc("/health", h.Health)

}

func InitRoutes_Proxy(router *http.ServeMux, h *models.ProxyHandlers) { // Internal only
	router.HandleFunc("/api/sdk/login", h.ProxyLogin)
	router.HandleFunc("/api/sdk/register", h.ProxyRegister)
	router.HandleFunc("/api/sdk/logout", h.ProxyLogout)
	router.HandleFunc("/api/sdk/token/refresh", h.ProxyRefreshToken)

}
