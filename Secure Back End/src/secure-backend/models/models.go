package models

type AuthConfigStruct struct {
	DATABASE_URL   string
	DB_DRIVER      string
	JWT_SECRET_KEY string
	PORT           string
	REDIS_ADDR     string
}

type DataConfigStruct struct {
	DATABASE_URL string
	DB_DRIVER    string
	PORT         string
}
