package models

type AuthConfigStruct struct {
	DATABASE_URL   string
	DB_DRIVER      string
	JWT_SECRET_KEY string
	PORT           string
	REDIS_ADDR     string
}
type ProxyConfigStruct struct {
	BASE_URL    string
	LOGIN_ROUTE string
}
