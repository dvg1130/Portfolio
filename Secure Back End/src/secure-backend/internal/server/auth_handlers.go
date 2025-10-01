package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dvg1130/Portfolio/secure-backend/internal/auth"
	"github.com/dvg1130/Portfolio/secure-backend/internal/helpers"
	validator "github.com/dvg1130/Portfolio/secure-backend/internal/validator/auth"
	"github.com/dvg1130/Portfolio/secure-backend/logs"
	"github.com/dvg1130/Portfolio/secure-backend/models"
	authdb "github.com/dvg1130/Portfolio/secure-backend/repo/auth_db"
)

// entry
func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("successful connection to server"))

}

// login
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {

	//decode json
	req, err := helpers.DecodeBody[models.Credentials](w, r)
	if err != nil {
		return
	}

	//check ip for lockout/blacklist

	ip := helpers.ClientIP(r)
	locked, _ := s.Redis.Exists(r.Context(), "lockout:"+ip).Result()
	if locked > 0 {
		logs.LogEvent(
			s.Logger, "warn", "Account lockout triggered", r,
			map[string]interface{}{
				"ip":         ip,
				"path":       r.URL.Path,
				"user_agent": r.UserAgent(),
			},
		)
		http.Error(w, "Too many failed attempts. Try again in 1 hour.", http.StatusTooManyRequests)
		return
	}

	// fetch user by username

	//query for user
	var storedHash, userRole, uuid string

	err = s.AUTH_DB.QueryRow(authdb.LoginUser, req.Username).Scan(&storedHash, &userRole, &uuid)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	//hashed pw check
	if !auth.CheckHashedPW(req.Password, storedHash) {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		logs.LogEvent(s.Logger, "warn", "Failed login attempt", r, map[string]interface{}{
			"username": req.Username,
		})
		// count failed attempt
		s.trackFailedAttempt(r.Context(), ip)
		return
	}

	//reset failed login counter
	s.Redis.Del(r.Context(), "fail:"+ip)

	//create jwt token
	w.Header().Set("Content-Type", "application/json")

	accesstoken, deviceid, err := auth.CreateAccessToken(req.Username, uuid, userRole)
	if err != nil {
		fmt.Println("error generating token", err)
		return
	}

	//create refresh token
	ctx := r.Context()

	refreshToken, exp, err := auth.CreateRefreshToken(req.Username, uuid, userRole)
	if err != nil {
		http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	//store refresh toekn and device_id in redis
	session := models.RefreshSession{
		RefreshToken: refreshToken,
		DeviceID:     deviceid,
		ExpiresAt:    exp,
	}

	sessionJSON, _ := json.Marshal(session)
	key := fmt.Sprintf("refresh:%s", req.Username)

	err = s.Redis.Set(ctx, key, sessionJSON, 7*24*time.Hour).Err()
	if err != nil {
		fmt.Println("failed to save refresh token", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	//set refresh token http only
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  exp,
		HttpOnly: true,
		Secure:   true,             // only over HTTPS
		Path:     "/token/refresh", // restrict usage
		SameSite: http.SameSiteStrictMode,
	})

	//sucessful login
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "login successful",
		"token":     accesstoken,
		"device_id": deviceid,
	})

}

// falied login tracker
func (s *Server) trackFailedAttempt(ctx context.Context, ip string) {
	key := "fail:" + ip
	count, _ := s.Redis.Incr(ctx, key).Result()
	if count == 1 {
		s.Redis.Expire(ctx, key, time.Hour)
	}
	if count >= 5 {
		// lockout for 1 hour and clear the fail counter
		s.Redis.Set(ctx, "lockout:"+ip, true, time.Hour)
		s.Redis.Del(ctx, key)
	}
}

// register
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {

	//decode body
	req, err := helpers.DecodeBody[models.Credentials](w, r)
	if err != nil {
		return
	}

	//len check
	if err := validator.AuthLenCheck(req.Username, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// existing username check
	exists, err := validator.ExistingUser(s.AUTH_DB, req.Username)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}

	//hash pw
	hashedPW, err := auth.HashPW(req.Password)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	//db insert
	_, err = s.AUTH_DB.Exec(authdb.RegisterUser, req.Username, hashedPW)
	if err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	//successfull registration
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("User registered successfully"))
}

// logout
func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// extract refresh token cookie
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		// delete from Redis
		s.Redis.Del(ctx, cookie.Value)

		// expire cookie on client
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   true,
			Path:     "/token/refresh",
			SameSite: http.SameSiteStrictMode,
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}

// refresh token
func (s *Server) TokenRefresh(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("successful connection to Logout"))
}
