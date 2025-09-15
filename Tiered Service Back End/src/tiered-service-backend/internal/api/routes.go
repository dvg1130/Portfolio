package api

import (
	"net/http"
	"tiered-service-backend/internal/middleware"
	"tiered-service-backend/internal/validators"

	"go.uber.org/zap"
)

type Handlers struct {
	Handler     http.HandlerFunc
	Login       http.HandlerFunc
	Register    http.HandlerFunc
	Logout      http.HandlerFunc
	Dashboard   http.HandlerFunc
	DashboardT1 http.HandlerFunc
	DashboardT2 http.HandlerFunc
	Admin       http.HandlerFunc
	Submit      http.HandlerFunc
	Refresh     http.HandlerFunc
	Health      http.HandlerFunc
	Logger      *zap.Logger
}

func InitRoutes(router *http.ServeMux, h *Handlers) {

	//routes
	router.HandleFunc("/", h.Handler)

	router.HandleFunc("/login", validators.Method(http.MethodPost, h.Login))

	router.HandleFunc("/register", validators.Method(http.MethodPost, h.Register))

	router.HandleFunc("/logout", validators.Method(http.MethodPost, h.Logout))

	//dashboard - basic(user)
	router.Handle("/dashboard",
		middleware.AuthMiddleware(
			validators.Method(http.MethodGet, http.HandlerFunc(h.Dashboard)),
		),
	)

	//dashboard - tier1
	router.Handle("/dashboard/t1",
		middleware.AuthMiddleware(
			middleware.RequireRole("tier1")(

				validators.Method(http.MethodGet, http.HandlerFunc(h.DashboardT1)),
			),
		),
	)

	//dashboard - tier2
	router.Handle("/dashboard/t2",
		middleware.AuthMiddleware(
			middleware.RequireRole("tier2")(

				validators.Method(http.MethodGet, http.HandlerFunc(h.DashboardT2)),
			),
		),
	)

	//admin
	router.Handle("/admin",
		middleware.AuthMiddleware(
			middleware.RequireRole("admin")(

				validators.Method(http.MethodGet, http.HandlerFunc(h.Admin)),
			),
		),
	)

	// submit
	router.Handle("/submit",
		middleware.LoggingMiddleware(h.Logger)(middleware.AuthMiddleware(
			middleware.PayloadLimiter(
				validators.Method(http.MethodPost, http.HandlerFunc(h.Submit)),
			),
		),
		),
	)

	//refresh token
	router.Handle("/token/refresh", validators.Method(http.MethodPost, http.HandlerFunc(h.Refresh)))

	//health
	router.Handle("/health", middleware.LoggingMiddleware(h.Logger)(http.HandlerFunc(h.Health)))

}
