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
	validator "github.com/dvg1130/Portfolio/secure-auth/internal/validator/auth"
	"github.com/dvg1130/Portfolio/secure-auth/logs"
	"github.com/dvg1130/Portfolio/secure-auth/models"
	authdb "github.com/dvg1130/Portfolio/secure-auth/repo/auth_db"
)

// entry
func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("successful connection to server"))

}

// login
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	//requestid
	// Use same type/key defined globally in your middleware package
	requestID, _ := r.Context().Value(middleware.RequestIDKey).(string)
	if requestID == "" {
		// Fallback to header if context is missing it
		requestID = r.Header.Get("X-Request-ID")
	}

	// Optional: fallback UUID to guarantee every request has an ID
	// if requestID == "" {
	//     requestID = uuid.New().String()
	// }

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

// register
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {

	type contextKey string

	const RequestIDKey contextKey = "requestID"
	claims := auth.GetClaimsFromContext(r.Context())
	requestID := r.Context().Value(RequestIDKey).(string)
	if claims != nil {
		requestID = claims["request_id"].(string)

	}

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
	logs.LogEvent(s.Logger, "info", "Successful Registration", r, map[string]interface{}{
		"path":       r.URL.Path,
		"user_agent": r.UserAgent(),
		"username":   req.Username,
		"category":   "auth",
		"action":     "user registeration",
		"request_id": requestID,
	})
}

// logout
func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {

	type LogoutRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	// get request id
	type contextKey string

	const RequestIDKey contextKey = "requestID"
	// var uuid string
	claim := auth.GetClaimsFromContext(r.Context())
	requestID := r.Context().Value(RequestIDKey).(string)
	uuid := r.Context().Value("uuid").(string)
	if claim != nil {
		requestID = claim["request_id"].(string)
		uuid = claim["uuid"].(string)
	}

	// verify/parse the refresh token to extract username
	claims, err := auth.VerifyToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	username, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "invalid token claims", http.StatusUnauthorized)
		return
	}
	// ctx := r.Context()
	ctx := auth.AddClaimsToContext(r.Context(), claims)
	r = r.WithContext(ctx)
	// fmt.Println(r)

	// delete session from Redis
	key := fmt.Sprintf("refresh:%s", username)
	deleted, _ := s.Redis.Del(r.Context(), key).Result()
	fmt.Println("Deleted keys:", deleted)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Refresh Token successfully"))
	logs.LogEvent(s.Logger, "info", "Successful Token Refesh", r, map[string]interface{}{
		"uuid":       uuid,
		"username":   username,
		"category":   "auth",
		"action":     "token refresh",
		"request_id": requestID,
	})
}

// refresh token
func (s *Server) TokenRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// extract refresh token from cookie
	cookie, err := r.Cookie("refresh_token")

	if err != nil {
		http.Error(w, "missing refresh token", http.StatusUnauthorized)
		return
	}

	fmt.Println("Incoming refresh token:", cookie.Value)

	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token"`
		DeviceID     string `json:"device_id"`
	}

	// Parse JSON body
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// get request id
	type contextKey string

	const RequestIDKey contextKey = "requestID"
	var uuid string
	claim := auth.GetClaimsFromContext(r.Context())
	requestID := r.Context().Value(RequestIDKey).(string)
	uuid = r.Context().Value("uuid").(string)
	if claim != nil {
		requestID = claim["request_id"].(string)
		uuid = claim["uuid"].(string)
	}

	// Verify refresh token (can be expired or validly signed)
	claims, err := auth.VerifyToken(req.RefreshToken)
	if err != nil && err.Error() != "expired token" {
		http.Error(w, "invalid token handler", http.StatusUnauthorized)
		fmt.Println(req.RefreshToken)
		fmt.Println(req.RefreshToken)
		return
	}

	username := claims["sub"].(string)
	key := fmt.Sprintf("refresh:%s", username)

	// Load stored session from Redis
	sessionJSON, err := s.Redis.Get(ctx, key).Result()
	if err != nil {
		http.Error(w, "no session found", http.StatusUnauthorized)
		return
	}

	var session models.RefreshSession
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		http.Error(w, "invalid session format", http.StatusUnauthorized)
		return
	}

	// Device ID check
	if session.DeviceID != req.DeviceID {
		http.Error(w, "device mismatch", http.StatusUnauthorized)
		return
	}

	// rotate logic
	newAccessToken, newRefreshToken, device_id, err := auth.RotateRefreshToken(ctx, s.Redis, req.RefreshToken)

	if err != nil {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// set new refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/token/refresh",
		SameSite: http.SameSiteStrictMode,
	})

	// return new access tokenls
	json.NewEncoder(w).Encode(map[string]string{

		"access_token": newAccessToken,
		"device_id":    device_id,
	})

	fmt.Println("new refresh token: ", newRefreshToken)
	logs.LogEvent(s.Logger, "info", "Successful token refresh", r, map[string]interface{}{
		"uuid":       uuid,
		"username":   username,
		"category":   "auth",
		"action":     "user logout",
		"request_id": requestID,
	})

}
