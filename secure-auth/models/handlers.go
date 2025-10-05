package models

import (
	"net/http"

	"go.uber.org/zap"
)

type AuthHandlers struct {
	//auth
	Handler      http.HandlerFunc
	Login        http.HandlerFunc
	Register     http.HandlerFunc
	Logout       http.HandlerFunc
	TokenRefresh http.HandlerFunc
	Logger       zap.SugaredLogger
}
