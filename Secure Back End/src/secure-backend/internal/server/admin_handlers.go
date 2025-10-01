package server

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/dvg1130/Portfolio/secure-backend/internal/auth"
	"github.com/dvg1130/Portfolio/secure-backend/internal/helpers"
	"github.com/dvg1130/Portfolio/secure-backend/logs"
	"github.com/dvg1130/Portfolio/secure-backend/models"
	authdb "github.com/dvg1130/Portfolio/secure-backend/repo/auth_db"
)

// GET  user lsit
func (s *Server) AdminGetAll(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}

	username := claims["username"].(string)
	role, ok := claims["role"].(string)
	if !ok || role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		logs.LogEvent(s.Logger, "warn", "Unauthorized access", r, map[string]interface{}{
			"username": username,
			"role":     role,
		})
		return
	}
	//query all users and roles
	rows, err := s.AUTH_DB.Query(authdb.AdminGetUsers)
	if err != nil {
		http.Error(w, "database query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Username, &u.Role); err != nil {
			http.Error(w, "row scan failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "encode failed: "+err.Error(), http.StatusInternalServerError)
	}
}

// GET one user
func (s *Server) AdminGetOne(w http.ResponseWriter, r *http.Request) {
	// ---- Claims / admin check ----
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	username := claims["username"].(string)
	role, ok := claims["role"].(string)
	if !ok || role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		logs.LogEvent(s.Logger, "warn", "Unauthorized access", r, map[string]interface{}{
			"username": username,
			"role":     role,
		})
		return
	}
	req, err := helpers.DecodeBody[models.User](w, r)
	if err != nil {
		return
	}

	// query for UUID of target user
	var uuid string
	err = s.AUTH_DB.QueryRow(authdb.AdminGetUuid, req.Username, req.Role).Scan(&uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "database query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// return user info
	user := models.User{
		Username: req.Username,
		Role:     req.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "encode failed: "+err.Error(), http.StatusInternalServerError)
	}

}

// UPDATE user role
func (s *Server) AdminUpdate(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	username := claims["username"].(string)
	role, ok := claims["role"].(string)
	if !ok || role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		logs.LogEvent(s.Logger, "warn", "Unauthorized access", r, map[string]interface{}{
			"username": username,
			"role":     role,
		})
		return
	}

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&models.RoleUpdate); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if models.RoleUpdate.Username == "" || models.RoleUpdate.OldRole == "" || models.RoleUpdate.NewRole == "" {
		http.Error(w, "username, old_role, and new_role are required", http.StatusBadRequest)
		return
	}

	// look up uuid , from username and old role
	var uuid string

	err := s.AUTH_DB.QueryRow(
		authdb.AdminGetUuid,
		models.RoleUpdate.Username, models.RoleUpdate.OldRole,
	).Scan(&uuid)
	if err == sql.ErrNoRows {
		http.Error(w, "user not found or old role mismatch", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "database query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//  Update the role by uuid
	_, err = s.AUTH_DB.Exec(
		authdb.AdminUpdateRole,
		models.RoleUpdate.NewRole, uuid,
	)
	if err != nil {
		http.Error(w, "failed to update role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"user role updated successfully"}`))

}
