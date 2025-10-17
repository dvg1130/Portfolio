package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dvg1130/Portfolio/secure-auth/internal/auth"
	"github.com/dvg1130/Portfolio/secure-auth/logs"
	"github.com/dvg1130/Portfolio/secure-auth/models"
)

// logout
func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {

	var req models.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var uuid string

	claims, err := auth.VerifyToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	// fmt.Println("claims: ", claims)

	username, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "invalid token claims", http.StatusUnauthorized)
		return
	}
	// ctx := r.Context()
	ctx := auth.AddClaimsToContext(r.Context(), claims)
	r = r.WithContext(ctx)
	// fmt.Println("claims: ", r)

	// delete session from Redis
	key := fmt.Sprintf("refresh:%s", username)
	deleted, _ := s.Redis.Del(r.Context(), key).Result()
	fmt.Println("Deleted keys:", deleted)

	w.WriteHeader(http.StatusOK)
	//successfull registration

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User logout successfully",
	})

	logs.LogEvent(s.Logger, "info", "Successful Logout", r, map[string]interface{}{
		"uuid":     uuid,
		"username": username,
		"category": "auth",
		"action":   "token refresh",
		// "request_id": requestID,
	})
}
