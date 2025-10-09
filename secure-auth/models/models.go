package models

type AuthConfigStruct struct {
	DATABASE_URL   string
	DB_DRIVER      string
	JWT_SECRET_KEY string
	PORT           string
	REDIS_ADDR     string
}
type ProxyConfigStruct struct {
	BASE_URL         string
	AUTH_LOGIN_ROUTE string
	LOGIN_ROUTE      string
	SDK_LOGIN_ROUTE  string
	REGUSTER_ROUTE   string
}
