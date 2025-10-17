package server

import (
	"encoding/json"
	"net/http"

	"github.com/dvg1130/Portfolio/secure-auth/internal/auth"
	"github.com/dvg1130/Portfolio/secure-auth/internal/helpers"
	validator "github.com/dvg1130/Portfolio/secure-auth/internal/validator/auth"
	"github.com/dvg1130/Portfolio/secure-auth/logs"
	"github.com/dvg1130/Portfolio/secure-auth/models"
	authdb "github.com/dvg1130/Portfolio/secure-auth/repo/auth_db"
)

// register
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {

	// type contextKey string

	// const RequestIDKey contextKey = "requestID"
	// claims := auth.GetClaimsFromContext(r.Context())
	// requestID := r.Context().Value(RequestIDKey).(string)
	// if claims != nil {
	// 	requestID = claims["request_id"].(string)

	// }

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
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
	})

	logs.LogEvent(s.Logger, "info", "Successful Registration", r, map[string]interface{}{
		"path":       r.URL.Path,
		"user_agent": r.UserAgent(),
		"username":   req.Username,
		"category":   "auth",
		"action":     "user registeration",
		// "request_id": requestID,
	})
}
