package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dvg1130/Portfolio/secure-auth/internal/auth"
	"github.com/dvg1130/Portfolio/secure-auth/internal/helpers"
	"github.com/dvg1130/Portfolio/secure-auth/internal/middleware"
	"github.com/dvg1130/Portfolio/secure-auth/logs"
	"github.com/dvg1130/Portfolio/secure-auth/models"
	authdb "github.com/dvg1130/Portfolio/secure-auth/repo/auth_db"
)

// login
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	//requestid

	requestID, _ := r.Context().Value(middleware.RequestIDKey).(string)
	if requestID == "" {
		// Fallback to header if context is missing it
		requestID = r.Header.Get("X-Request-ID")
	}

	claims := auth.GetClaimsFromContext(r.Context())
	if claims != nil {
		if rid, ok := claims["request_id"].(string); ok && rid != "" {
			requestID = rid
		}
	}
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
				"category":   "autht",
				"action":     "account lockout",
				"request_id": requestID,
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
			"uuid":       uuid,
			"username":   req.Username,
			"category":   "auth",
			"action":     "failed login attempt",
			"request_id": requestID,
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
	fmt.Println("access token from auth_handler", accesstoken)
	fmt.Println("refresh token from auth_handler", refreshToken)

	//store refresh token and device_id in redis
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
		Secure:   true,                      // only over HTTPS
		Path:     "/api/auth/token/refresh", // restrict usage
		SameSite: http.SameSiteStrictMode,
	})

	//sucessful login
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "login successful",
		"token":         accesstoken,
		"device_id":     deviceid,
		"refresh_token": refreshToken,
		"expires":       exp,
	})

	logs.LogEvent(s.Logger, "info", "Successful Login", r, map[string]interface{}{
		"uuid":       uuid,
		"username":   req.Username,
		"category":   "auth",
		"action":     "user login",
		"request_id": requestID,
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
