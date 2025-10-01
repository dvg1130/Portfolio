package server

import (
	"database/sql"
	"net/http"

	"github.com/dvg1130/Portfolio/secure-backend/internal/api"
	"github.com/dvg1130/Portfolio/secure-backend/internal/helpers"
	"github.com/dvg1130/Portfolio/secure-backend/internal/middleware"
	"github.com/dvg1130/Portfolio/secure-backend/models"
	redisdb "github.com/dvg1130/Portfolio/secure-backend/repo/redis_db"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Server struct {
	Router  *http.ServeMux
	AUTH_DB *sql.DB
	Data_DB *sql.DB
	Redis   *redis.Client
	Logger  *zap.SugaredLogger
	S       *sql.DB
}

func AppServer(auth_db *sql.DB, data_db *sql.DB, logger *zap.SugaredLogger) *Server {
	rdb := redisdb.RedisClient()
	s := &Server{
		Router:  http.NewServeMux(),
		AUTH_DB: auth_db,
		Data_DB: data_db,
		Redis:   rdb,
		Logger:  logger,
	}

	api.InitRoutes_Auth(s.Router, &models.AuthHandlers{
		Handler:  s.Handler,
		Login:    s.Login,
		Register: s.Register,
		Logout:   s.Logout,
	})

	api.InitRoutes_Data(s.Router, s.Data_DB, &models.DataHandlers{
		//snakes
		SnakeGetAll: s.SnakeGetAll,
		SnakeGetOne: s.SnakeGetOne,
		SnakePost:   s.SnakePost,
		SnakeUpdate: s.SnakeUpdate,
		SnakeDelete: s.SnakeDelete,

		//feeds
		SnakeFeedGet:    s.SnakeFeedGet,
		SnakeFeedPost:   s.SnakeFeedPost,
		SnakeFeedUpdate: s.SnakeFeedUpdate,
		SnakeFeedDelete: s.SnakeFeedDelete,

		//health
		SnakeHealthGet:    s.SnakeHealthGet,
		SnakeHealthPost:   s.SnakeHealthPost,
		SnakeHealthUpdate: s.SnakeHealthUpdate,
		SnakeHealthDelete: s.SnakeHealthDelete,

		S: s.Data_DB,
	})

	api.InitRoutes_Admin(s.Router, &models.AdminHandlers{
		AdminGetAll: s.AdminGetAll,
		AdminGetOne: s.AdminGetOne,
		AdminUpdate: s.AdminUpdate,
	})

	api.InitRoutes_Breeding(s.Router, &models.BreedingHandlers{
		//breeding
		SnakeBreedGetBySnake: s.SnakeBreedGetBySnake,
		SnakeBreedGetAll:     s.SnakeBreedGetAll,
		SnakeBreedGetOne:     s.SnakeBreedGetOne,
		SnakeBreedPost:       s.SnakeBreedPost,
		SnakeBreedUpdate:     s.SnakeBreedUpdate,
		SnakeBreedDelete:     s.SnakeBreedDelete,
	})

	s.Router = helpers.ServeMuxWrapper(
		s.Router,
		middleware.SecurityHeaders,
		middleware.LoggingMiddleware(logger),
	)

	return s
}
