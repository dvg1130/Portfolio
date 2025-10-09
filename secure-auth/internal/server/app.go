package server

import (
	"database/sql"
	"net/http"

	"github.com/dvg1130/Portfolio/secure-auth/authsdk"
	"github.com/dvg1130/Portfolio/secure-auth/internal/api"
	"github.com/dvg1130/Portfolio/secure-auth/internal/helpers"
	"github.com/dvg1130/Portfolio/secure-auth/internal/middleware"
	"github.com/dvg1130/Portfolio/secure-auth/models"
	redisdb "github.com/dvg1130/Portfolio/secure-auth/repo/redis_db"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Server struct {
	Router      *http.ServeMux
	AUTH_DB     *sql.DB
	AUTH_CLIENT *authsdk.SDKClient
	Redis       *redis.Client
	Logger      *zap.SugaredLogger
	S           *sql.DB
}

func AppServer(auth_db *sql.DB, logger *zap.SugaredLogger) *Server {
	rdb := redisdb.RedisClient()
	s := &Server{
		Router:  http.NewServeMux(),
		AUTH_DB: auth_db,
		Redis:   rdb,
		Logger:  logger,
	}

	api.InitRoutes_Proxy(s.Router, &models.AuthHandlers{
		ProxyLogin: s.ProxyLogin,
	})

	api.InitRoutes_Auth(s.Router, &models.AuthHandlers{
		Handler: s.Handler,

		Login:        s.Login,
		Register:     s.Register,
		Logout:       s.Logout,
		TokenRefresh: s.TokenRefresh,
		Logger:       *s.Logger,
		Health:       s.Health,
	})

	s.Router = helpers.ServeMuxWrapper(
		s.Router,
		middleware.SecurityHeaders,
		middleware.RequestIDMiddleware(s.Logger),
		middleware.IncomingLoggerMiddleware(s.Logger),
	)

	return s
}
