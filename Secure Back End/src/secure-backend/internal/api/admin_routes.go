package api

import (
	"net/http"

	"github.com/dvg1130/Portfolio/secure-backend/internal/middleware"
	validator "github.com/dvg1130/Portfolio/secure-backend/internal/validator/util"
	"github.com/dvg1130/Portfolio/secure-backend/models"
)

func InitRoutes_Admin(router *http.ServeMux, h *models.AdminHandlers) {

	//routes

	router.Handle("/admin/users/all",
		middleware.AuthMiddleware(
			middleware.RequireRole("admin")(
				validator.Method(http.MethodGet,
					http.HandlerFunc(h.AdminGetAll)),
			),
		),
	)

	router.Handle("/admin/user/one",
		middleware.AuthMiddleware(
			middleware.RequireRole("admin")(
				validator.Method(http.MethodGet,
					http.HandlerFunc(h.AdminGetOne)),
			),
		),
	)
	router.Handle("/admin/user/update",
		middleware.AuthMiddleware(
			middleware.RequireRole("admin")(
				validator.Method(http.MethodPatch,
					http.HandlerFunc(h.AdminUpdate)),
			),
		),
	)
}
