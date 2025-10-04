package api

import (
	"net/http"

	"github.com/dvg1130/Portfolio/secure-backend/internal/middleware"
	validator "github.com/dvg1130/Portfolio/secure-backend/internal/validator/util"
	"github.com/dvg1130/Portfolio/secure-backend/models"
	"go.uber.org/zap"
)

func InitRoutes_Admin(router *http.ServeMux, z *zap.SugaredLogger, h *models.AdminHandlers) {

	//routes

	router.Handle("/admin/users/all",
		middleware.AuthMiddleware(
			middleware.RequireRole("admin", z)(
				validator.Method(http.MethodGet,
					http.HandlerFunc(h.AdminGetAll)),
			),
		),
	)

	router.Handle("/admin/user/one",
		middleware.AuthMiddleware(
			middleware.RequireRole("admin", z)(
				validator.Method(http.MethodGet,
					http.HandlerFunc(h.AdminGetOne)),
			),
		),
	)
	router.Handle("/admin/user/update",
		middleware.AuthMiddleware(
			middleware.RequireRole("admin", z)(
				validator.Method(http.MethodPatch,
					http.HandlerFunc(h.AdminUpdate)),
			),
		),
	)
}
