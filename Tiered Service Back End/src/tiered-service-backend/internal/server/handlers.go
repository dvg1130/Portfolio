package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"tiered-service-backend/internal/auth"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

// creds struct
type Credentials struct {
	Username  string `json:"username"`
	Passsword string `json:"password"`
}

// login struct
type Login struct {
	HashedPW string
	JWTToken string
}

// handler func
func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sucessful connection to server!"))

}

// login handler
func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Login route!"))
	var cred Credentials

	//decode json
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// check ip for lockout
	ip := clientIP(r)
	locked, _ := s.Redis.Exists(r.Context(), "lockout:"+ip).Result()
	if locked > 0 {
		http.Error(w, "Too many failed attempts. Try again in 1 hour.", http.StatusTooManyRequests)
		return
	}

	// fetch user by username
	var storedHash, userRole string
	type ctxKey string
	const loggerKey ctxKey = "logger"
	logger, ok := r.Context().Value(loggerKey).(*zap.SugaredLogger)
	if !ok {
		// fallback if logger not found
		logger = zap.NewExample().Sugar()
	}

	err := s.DB.QueryRow("SELECT password, role FROM users WHERE username = ?", cred.Username).Scan(&storedHash, &userRole)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warnw("Failed login attempt",
				"timestamp", time.Now().Format(time.RFC3339),
				"username", cred.Username,
				"ip", r.RemoteAddr,
				"path", r.URL.Path,
			)

			// count failed attempt
			s.trackFailedAttempt(r.Context(), ip)
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// compare provided password with stored hash
	if !auth.CheckHashedPW(cred.Passsword, storedHash) {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	// reset failure counter on successful login
	s.Redis.Del(r.Context(), "fail:"+ip)

	// jwt token

	w.Header().Set("Content-Type", "application/json")

	token, err := auth.CreateToken(cred.Username, userRole)
	if err != nil {
		fmt.Println("error generating token", err)
		return
	}

	//create refresh token

	ctx := r.Context()
	refreshToken, err := auth.CreateRefreshToken(cred.Username, userRole)
	if err != nil {
		http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	//store refresh token redis
	err = s.Redis.Set(ctx, cred.Username, refreshToken, 7*24*time.Hour).Err()
	if err != nil {
		fmt.Println("failed to save refresh token:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}
	//set refresh token as http-only
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,             // only over HTTPS
		Path:     "/token/refresh", // restrict usage
		SameSite: http.SameSiteStrictMode,
	})

	// success
	json.NewEncoder(w).Encode(map[string]string{
		"message": "login successful",
		"token":   token,
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

// extract client IP
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

// register handler
func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Register route!"))
	var cred Credentials

	//decode json
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	//length check
	if len(cred.Username) < 8 || len(cred.Passsword) < 8 {
		err := http.StatusNotAcceptable
		http.Error(w, "Invalid Username/password", err)
		return
	}

	//existing user check
	var exists bool
	err := s.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", cred.Username).Scan(&exists)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	//hash pw
	hashedPW, err := auth.HashPW(cred.Passsword)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	// insert into DB (using prepared statement)
	_, err = s.DB.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		cred.Username,
		hashedPW,
	)
	if err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	//success registeration message
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))

}

// logout
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
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

// dashboard free tier
func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return

	}
	username := claims["username"].(string)
	role := claims["role"].(string)

	fmt.Fprintf(w, "Welcome %s! Your role is %s", username, role)

}

// dashboard tier1
func (s *Server) dashboardT1(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return

	}
	username := claims["username"].(string)
	role := claims["role"].(string)

	fmt.Fprintf(w, "Welcome %s to Tier 1! Your role is %s", username, role)

}

//dashboard tier 2

func (s *Server) dashboardT2(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return

	}
	username := claims["username"].(string)
	role := claims["role"].(string)

	fmt.Fprintf(w, "Welcome %s to Tier 2! Your role is %s", username, role)

}

// admin panel
func (s *Server) admin(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return

	}
	username := claims["username"].(string)
	role := claims["role"].(string)

	fmt.Fprintf(w, "Welcome %s to Admin Panel! Your role is %s", username, role)

}

// submit
func (s *Server) submit(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return

	}
	username := claims["username"].(string)
	// role := claims["role"].(string)

	fmt.Fprintf(w, "%s, your upload was successful", username)

}

// refresh token
func (s *Server) refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// extract refresh token from cookie
	cookie, err := r.Cookie("refresh_token")

	fmt.Println("Incoming refresh token:", cookie.Value)
	if err != nil {
		http.Error(w, "missing refresh token", http.StatusUnauthorized)
		return
	}

	// rotate logic
	newAccessToken, newRefreshToken, err := auth.RotateRefresh(ctx, s.Redis, cookie.Value)
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

	// return new access token
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}

// health
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
