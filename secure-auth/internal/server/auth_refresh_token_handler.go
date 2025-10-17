package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dvg1130/Portfolio/secure-auth/internal/auth"
	"github.com/dvg1130/Portfolio/secure-auth/logs"
	"github.com/dvg1130/Portfolio/secure-auth/models"
)

// refresh token
func (s *Server) TokenRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse JSON body
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Verify refresh token (can be expired or validly signed)
	claims, err := auth.VerifyToken(req.RefreshToken)
	if err != nil && err.Error() != "expired token" {
		http.Error(w, "invalid token handler", http.StatusUnauthorized)
		fmt.Println(req.RefreshToken)

		return
	}
	uuid := claims["uuid"].(string)
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
	exp := time.Now().Add(7 * 24 * time.Hour)
	// set new refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "refresh_token",
		Value: newRefreshToken,
		// Expires:  time.Now().Add(7 * 24 * time.Hour),
		Expires:  exp,
		HttpOnly: true,
		Secure:   true,
		Path:     "/api/auth/token/refresh",
		SameSite: http.SameSiteStrictMode,
	})

	// return new atokens
	json.NewEncoder(w).Encode(map[string]interface{}{

		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
		"device_id":     device_id,
		"expires":       exp,
	})

	fmt.Println("new refresh token: ", newRefreshToken)
	logs.LogEvent(s.Logger, "info", "Successful token refresh", r, map[string]interface{}{
		"uuid":     uuid,
		"username": username,
		"category": "auth",
		"action":   "user logout",
		// "request_id": requestID,
	})

}
