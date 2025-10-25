package server

import (
	"database/sql"
	"net/http"

	"github.com/dvg1130/Portfolio/secure-auth/authsdk"
	"github.com/dvg1130/Portfolio/secure-backend/internal/api"
	"github.com/dvg1130/Portfolio/secure-backend/internal/helpers"
	"github.com/dvg1130/Portfolio/secure-backend/internal/middleware"
	"github.com/dvg1130/Portfolio/secure-backend/models"
	"go.uber.org/zap"
)

type Server struct {
	Router *http.ServeMux

	AuthClient *authsdk.SDKClient
	Data_DB    *sql.DB

	Logger *zap.SugaredLogger
	S      *sql.DB
}

func AppServer(authClient *authsdk.SDKClient, data_db *sql.DB, logger *zap.SugaredLogger) *Server {

	s := &Server{
		Router:     http.NewServeMux(),
		AuthClient: authClient,
		Data_DB:    data_db,

		Logger: logger,
	}

	// api.InitRoutes_Auth(s.Router, &models.AuthHandlers{
	// 	Handler:      s.Handler,
	// 	Login:        s.Login,
	// 	Register:     s.Register,
	// 	Logout:       s.Logout,
	// 	TokenRefresh: s.TokenRefresh,
	// })

	api.InitRoutes_Data(s.Router, s.Data_DB, s.Logger, &models.DataHandlers{
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
	},
	)

	api.InitRoutes_Auth(s.Router, &models.AuthHandlers{
		RegisterHandler: s.RegisterHandler,
		LoginHandler:    s.LoginHandler,
		LogoutHandler:   s.LogoutHandler,
		TokenRefresh:    s.TokenRefresh,
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
		middleware.IncomingLoggerMiddleware(s.Logger),
	)

	return s
}
