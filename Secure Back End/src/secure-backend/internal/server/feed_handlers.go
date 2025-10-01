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

// GET snake feed recs
func (s *Server) SnakeFeedGet(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	userUUID := claims["uuid"].(string)

	// Expect {"sid":"<snake-id>"} in request body
	defer r.Body.Close()
	req, err := helpers.DecodeBody[models.SnakeSid](w, r)

	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	// First, get suid for this user’s snake
	var suid string
	err = s.Data_DB.QueryRow(
		datadb.GetSuid,
		userUUID, req.Sid,
	).Scan(&suid)
	if err == sql.ErrNoRows {
		http.Error(w, "snake not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error retrieving snake id: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Now get all feeding records for that snake
	rows, err := s.Data_DB.Query(
		datadb.GetSnakeFeeds,
		suid,
	)
	if err != nil {
		http.Error(w, "error retrieving feeding records: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var feeds []models.SnakeFeed
	for rows.Next() {
		var f models.SnakeFeed
		if err := rows.Scan(&f.Sid,
			&f.FeedDate, &f.PreyType, &f.PreySize, &f.Notes); err != nil {
			http.Error(w, "error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		feeds = append(feeds, f)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "row iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(feeds); err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
	}

}

// POST snake feed rec
func (s *Server) SnakeFeedPost(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	userUUID := claims["uuid"].(string)

	req, err := helpers.DecodeBody[models.SnakeFeed](w, r)
	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	//find snake suid
	var suid string
	err = s.Data_DB.QueryRow(
		datadb.GetSuid,
		userUUID, req.Sid,
	).Scan(&suid)
	if err == sql.ErrNoRows {
		http.Error(w, "snake not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error finding snake: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//  insert feeding record with suid as foreign key
	_, err = s.Data_DB.Exec(
		datadb.AddSnakeFeed,
		suid, req.Sid, req.FeedDate, req.PreyType, req.PreySize, req.Notes,
	)
	if err != nil {
		http.Error(w, "error inserting feeding record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"feeding record added"}`))
}

//UPDATE snake feed rec

func (s *Server) SnakeFeedUpdate(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	userUUID := claims["uuid"].(string)

	defer r.Body.Close()

	req, err := helpers.DecodeBody[models.UpdateSnakeFeed](w, r)
	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate: need sid + feed_date + at least one update field
	if req.Sid == "" || req.FeedDate == "" {
		http.Error(w, "sid and feed_date are required", http.StatusBadRequest)
		return
	}
	if req.PreyType == nil && req.PreySize == nil && req.Notes == nil {
		http.Error(w, "at least one field to update must be provided", http.StatusBadRequest)
		return
	}

	// get suid from snakes
	var suid string
	err = s.Data_DB.QueryRow(
		datadb.GetSuid,
		userUUID, req.Sid,
	).Scan(&suid)
	if err == sql.ErrNoRows {
		http.Error(w, "snake not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error retrieving snake: "+err.Error(), http.StatusInternalServerError)
		return
	}

	updates := []string{}
	args := []interface{}{}

	if req.PreyType != nil {
		updates = append(updates, "prey_type = ?")
		args = append(args, *req.PreyType)
	}
	if req.PreySize != nil {
		updates = append(updates, "prey_size = ?")
		args = append(args, *req.PreySize)
	}
	if req.Notes != nil {
		updates = append(updates, "notes = ?")
		args = append(args, *req.Notes)
	}

	// add WHERE args (suid + feed_date)
	args = append(args, suid, req.FeedDate)

	query := fmt.Sprintf(
		datadb.UpdateSnakeFeed,
		strings.Join(updates, ", "),
	)

	//  Execute update
	result, err := s.Data_DB.Exec(query, args...)
	if err != nil {
		http.Error(w, "error updating feeding record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "could not confirm update: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "feeding record not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"feeding record updated"}`))

}

// DELETE snake feed rec
func (s *Server) SnakeFeedDelete(w http.ResponseWriter, r *http.Request) {
	// ---- 1. Verify JWT claims ----
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	userUUID := claims["uuid"].(string)

	// ---- 2. Parse and validate request body ----
	var req models.SnakeFeedDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Sid == "" || req.Date == "" {
		http.Error(w, "sid and date are required", http.StatusBadRequest)
		return
	}

	// ---- 3. Get suid for this snake ----
	var suid string
	err := s.Data_DB.QueryRow(
		datadb.GetSuid,
		userUUID, req.Sid,
	).Scan(&suid)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "snake not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error retrieving snake", http.StatusInternalServerError)
		return

	}

	// ---- 4. Delete feeding record ----
	res, err := s.Data_DB.Exec(
		datadb.DeleteSnakeFeed,
		suid, req.Date,
	)

	if err != nil {
		http.Error(w, "error deleting feeding record", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "feeding record not found", http.StatusNotFound)
		return
	}

	// ---- 5. Respond success ----
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"feeding record deleted successfully"}`))

}
