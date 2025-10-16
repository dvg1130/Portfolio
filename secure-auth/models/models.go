package models

type AuthConfigStruct struct {
	DATABASE_URL   string
	DB_DRIVER      string
	JWT_SECRET_KEY string
	PORT           string
	REDIS_ADDR     string
}
type ProxyConfigStruct struct {
	BASE_URL string

	AUTH_LOGIN_ROUTE string
	SDK_LOGIN_ROUTE  string

	AUTH_REGISTER_ROUTE string
	SDK_REGISTER_ROUTE  string

	AUTH_LOGOUT_ROUTE string
	SDK_LOGOUT_ROUTE  string
}
