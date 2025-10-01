package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dvg1130/Portfolio/secure-backend/internal/auth"
	"github.com/dvg1130/Portfolio/secure-backend/internal/helpers"
	"github.com/dvg1130/Portfolio/secure-backend/models"
	datadb "github.com/dvg1130/Portfolio/secure-backend/repo/data_db"
)

// GET all snakes
func (s *Server) SnakeGetAll(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}

	uuid := claims["uuid"].(string)

	// get all snakes for this owner
	rows, err := s.Data_DB.Query(
		datadb.GetSnakes, uuid,
	)
	if err != nil {
		http.Error(w, "error retrieving snakes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var snakeList []models.SnakesListItem
	for rows.Next() {
		var sItem models.SnakesListItem
		if err := rows.Scan(&sItem.Sid, &sItem.Species); err != nil {
			http.Error(w, "error scanning row", http.StatusInternalServerError)
			return
		}
		snakeList = append(snakeList, sItem)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "error iterating rows", http.StatusInternalServerError)
		return
	}
	// Return JSON list
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(snakeList); err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

// GET one snake
func (s *Server) SnakeGetOne(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	uuid := claims["uuid"].(string)

	defer r.Body.Close()
	var req models.Snake
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// query a single snake
	var snake models.Snake
	err := s.Data_DB.QueryRow(
		datadb.GetSnake,
		req.SnakeId, uuid,
	).Scan(
		&snake.SnakeId, &snake.Species, &snake.Sex,
		&snake.Age, &snake.Genes, &snake.Notes,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "snake not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error retrieving snake: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(snake); err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
	}
}

// POST snake
func (s *Server) SnakePost(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return

	}
	uuid := claims["uuid"].(string)
	// role := claims["role"].(string)

	//decode json
	var req models.Snake
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body"+err.Error(), http.StatusBadRequest)
		return
	}

	// existing snake check (unique sid per user)
	var exists bool
	err := s.Data_DB.QueryRow(
		datadb.SnakeExists,
		req.SnakeId, uuid,
	).Scan(&exists)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Snake ID already exists for this user", http.StatusConflict)
		return
	}

	//rbac middleware will wrap post method at router level, additional role vs max allowed check.. if then error, else then return

	// insert into DB (using prepared statement)
	_, err = s.Data_DB.Exec(
		//maybe add snakeuuid for unique snake id across all users to prevent data leak
		datadb.AddSnake,
		uuid,
		req.SnakeId,
		req.Species,
		req.Sex,
		req.Age,
		req.Genes,
		req.Notes,
	)

	if err != nil {
		http.Error(w, "error creating new snake", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Upload Successful")

}

// UPDATE snake
func (s *Server) SnakeUpdate(w http.ResponseWriter, r *http.Request) {
	//claims
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	uuid := claims["uuid"].(string)

	//decode json
	defer r.Body.Close()

	var req models.UpdateSnake

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Sid == "" {
		http.Error(w, "sid is required", http.StatusBadRequest)
		return
	}

	// build  clause/case
	sets := []string{}
	args := []any{}

	if req.Species != nil {
		sets = append(sets, "species = ?")
		args = append(args, *req.Species)
	}
	if req.Sex != nil {
		sets = append(sets, "sex = ?")
		args = append(args, *req.Sex)
	}
	if req.Age != nil {
		sets = append(sets, "age = ?")
		args = append(args, *req.Age)
	}
	if req.Genes != nil {
		sets = append(sets, "genes = ?")
		args = append(args, *req.Genes)
	}
	if req.Notes != nil {
		sets = append(sets, "notes = ?")
		args = append(args, *req.Notes)
	}

	if len(sets) == 0 {
		http.Error(w, "no fields to update", http.StatusBadRequest)
		return
	}

	// add WHERE arguments: sid and owner_uuid
	args = append(args, req.Sid, uuid)

	query := fmt.Sprintf(datadb.UpdateSnake,
		strings.Join(sets, ", "))

	if _, err := s.Data_DB.Exec(query, args...); err != nil {
		http.Error(w, "update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"updated"}`))
}

// DELETE snake
func (s *Server) SnakeDelete(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	uuid := claims["uuid"].(string)

	defer r.Body.Close()

	req, err := helpers.DecodeBody[models.SnakeSid](w, r)

	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Sid == "" {
		http.Error(w, "sid is required", http.StatusBadRequest)
		return
	}

	result, err := s.Data_DB.Exec(
		datadb.DeleteSnake,
		req.Sid, uuid,
	)
	if err != nil {
		http.Error(w, "error deleting snake: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "could not confirm deletion: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "snake not found or not owned by user", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 – success, no response body
	w.Write([]byte("snakes deleted"))
}
