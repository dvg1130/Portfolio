package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dvg1130/Portfolio/secure-auth/models"

	"github.com/dvg1130/Portfolio/secure-backend/internal/helpers"
	// "github.com/dvg1130/Portfolio/secure-backend/models"
)

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {

	// 	//decode json
	req, err := helpers.DecodeBody[models.Credentials](w, r)
	if err != nil {
		return
	}
	var creds = models.Credentials{
		Username: req.Username,
		Password: req.Password,
	}

	loginResp, err := s.AuthClient.SDKLogin(creds)
	if err != nil {
		fmt.Println("Error from snake app Login Handler")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Log or store the device ID
	fmt.Println("Device ID:", loginResp.DeviceID)

	// Optionally set it as cookie or header
	w.Header().Set("X-Device-ID", loginResp.DeviceID)
	//set refresh token http only
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    loginResp.RefreshToken,
		Expires:  loginResp.Expires,
		HttpOnly: true,
		Secure:   true,                      // only over HTTPS
		Path:     "/api/auth/token/refresh", // restrict usage
		SameSite: http.SameSiteStrictMode,
	})

	// Respond to client
	json.NewEncoder(w).Encode(loginResp)
	fmt.Println("login respones:", loginResp)
}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	//extract refresh token from cookie
	cookie := "refresh_token"
	refreshtokenString := helpers.ExtractCookieValue(cookie, w, r)
	fmt.Println("refresh token string: ", refreshtokenString)

	//extract access token from header
	accesstokenString := helpers.ExtractHeader_Token(w, r)
	fmt.Println("access token string: ", refreshtokenString)

	var logoutRequest = models.LogoutRequest{
		AccessToken:  accesstokenString,
		RefreshToken: refreshtokenString,
	}

	fmt.Println(logoutRequest.RefreshToken)

	logoutResp, err := s.AuthClient.SDKLogout(logoutRequest)
	if err != nil {
		fmt.Println("Error from snake app refresh tokenHandler", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//set refresh token http only
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,                      // only over HTTPS
		Path:     "/api/auth/token/refresh", // restrict usage
		SameSite: http.SameSiteStrictMode,
	})

	// Respond to client
	json.NewEncoder(w).Encode(logoutResp)
	fmt.Println("logout response:", logoutResp)

}

func (s *Server) TokenRefresh(w http.ResponseWriter, r *http.Request) {

	// get token from auth header
	tokenString := helpers.ExtractHeader_Token(w, r)

	req, err := helpers.DecodeBody[models.RefreshRequest](w, r)
	if err != nil {
		return
	}

	var refreshRequest = models.RefreshRequest{
		AccessToken:  tokenString,
		RefreshToken: req.RefreshToken,
		DeviceID:     req.DeviceID,
	}
	fmt.Println("refresh request :", refreshRequest)

	refreshResp, err := s.AuthClient.SDKRefreshToken(refreshRequest)
	if err != nil {
		fmt.Println("Error from snake app refresh tokenHandler", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Log or store the device ID
	fmt.Println("Device ID:", refreshResp.DeviceID)

	// Optionally set it as cookie or header
	w.Header().Set("X-Device-ID", refreshResp.DeviceID)
	//set refresh token http only
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshResp.RefreshToken,
		Expires:  refreshResp.Expires,
		HttpOnly: true,
		Secure:   true,                      // only over HTTPS
		Path:     "/api/auth/token/refresh", // restrict usage
		SameSite: http.SameSiteStrictMode,
	})

	// Respond to client
	json.NewEncoder(w).Encode(refreshResp)
	fmt.Println("refresh respones:", refreshResp)

}

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// 	//decode json
	req, err := helpers.DecodeBody[models.Credentials](w, r)
	if err != nil {
		return
	}
	var creds = models.Credentials{
		Username: req.Username,
		Password: req.Password,
	}
	registerResp, err := s.AuthClient.SDKRegister(creds)
	if err != nil {
		fmt.Println("Error from snake app register Handler")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Respond to client
	json.NewEncoder(w).Encode(registerResp)
	fmt.Println("register respones:", registerResp)

}
