package api

import (
	"net/http"

	validator "github.com/dvg1130/Portfolio/secure-backend/internal/validator/util"
	"github.com/dvg1130/Portfolio/secure-backend/models"
)

func InitRoutes_Auth(router *http.ServeMux, h *models.AuthHandlers) {

	//routes
	router.HandleFunc("/", h.Handler)
	router.HandleFunc("/login", validator.Method(http.MethodPost, h.Login))
	router.HandleFunc("/register", validator.Method(http.MethodPost, h.Register))
	router.HandleFunc("/logout", validator.Method(http.MethodPost, h.Logout))
}
