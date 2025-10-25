package models

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Login struct {
	HashedPW string
	JWTToken string
}

type RefreshSession struct {
	RefreshToken string `json:"refresh_token"`
	DeviceID     string `json:"device_id"`
}

type User struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

var RoleUpdate struct {
	Username string `json:"username"`
	OldRole  string `json:"old_role"`
	NewRole  string `json:"new_role"`
}
