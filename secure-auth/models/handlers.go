package models

import (
	"net/http"

	"go.uber.org/zap"
)

type AuthHandlers struct {
	//auth
	Handler      http.HandlerFunc
	Login        http.HandlerFunc
	Health       http.HandlerFunc
	Register     http.HandlerFunc
	Logout       http.HandlerFunc
	TokenRefresh http.HandlerFunc
	Logger       zap.SugaredLogger
}
type ProxyHandlers struct {
	//auth

	ProxyLogin    http.HandlerFunc
	ProxyRegister http.HandlerFunc
	ProxyLogout   http.HandlerFunc
}
