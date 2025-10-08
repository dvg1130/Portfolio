package models

import (
	"net/http"
	"time"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Login struct {
	HashedPW string
	JWTToken string
}

type RefreshSession struct {
	RefreshToken string    `json:"refresh_token"`
	DeviceID     string    `json:"device_id"`
	ExpiresAt    time.Time `json:"expires_at"`
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

// type LoginResponse struct {
// 	AccessToken  string `json:"access_token"`
// 	RefreshToken string `json:"refresh_token"`
// 	ExpiresIn    int64  `json:"expires_in"`
// 	Username       string `json:"username"`
// 	DeviceID     string `json:"device_id"`

// }

type LoginResponse struct {
	Username     string `json:"username"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	DeviceID     string `json:"device_id,omitempty"`
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}
