package server

import (
	"database/sql"
	"net/http"
	"tiered-service-backend/internal/api"
	"tiered-service-backend/internal/middleware"
	"tiered-service-backend/repository/redisdb"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Server struct {
	Router *http.ServeMux
	DB     *sql.DB
	Redis  *redis.Client
	Logger *zap.Logger
}

func NewServer(db *sql.DB, logger *zap.Logger) *Server {
	rdb := redisdb.NewClient()

	s := &Server{
		Router: http.NewServeMux(),
		//db
		DB:     db,
		Redis:  rdb,
		Logger: logger,
	}
	api.InitRoutes(s.Router, &api.Handlers{
		Handler:     s.handler,
		Login:       s.login,
		Register:    s.register,
		Logout:      s.logout,
		Dashboard:   s.dashboard,
		DashboardT1: s.dashboardT1,
		DashboardT2: s.dashboardT2,
		Admin:       s.admin,
		Submit:      s.submit,
		Refresh:     s.refresh,
		Health:      s.health,
		Logger:      s.Logger,
	})

	s.Router = ServeMuxWrapper(
		s.Router,
		middleware.SecurityHeaders,
		middleware.LoggingMiddleware(logger),
	)

	return s
}
