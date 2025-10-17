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

type LoginResponse struct {
	Username     string    `json:"username"`
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	DeviceID     string    `json:"device_id,omitempty"`
	Expires      time.Time `json:"expires,omitempty"`
}

type RegisterResponse struct {
	Message string `json:"messsage,omitempty"`
}

type LogoutRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type LogoutResponse struct {
	Message string `json:"messsage,omitempty"`
}

type RefreshRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	DeviceID     string `json:"device_id"`
}
type RefreshResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	DeviceID     string    `json:"device_id"`
	Expires      time.Time `json:"expires"`
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}
