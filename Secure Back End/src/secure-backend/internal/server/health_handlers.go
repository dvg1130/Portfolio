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

// GET snake Health recs
func (s *Server) SnakeHealthGet(w http.ResponseWriter, r *http.Request) {
	//get uuid from claims
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	userUUID := claims["uuid"].(string)

	//decode json
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

	// get suid
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

	//query health records for suid
	rows, err := s.Data_DB.Query(
		datadb.GetSnakeHealth,
		suid,
	)
	if err != nil {
		http.Error(w, "error retrieving health records: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	//get health record
	var records []models.HealthRecord

	for rows.Next() {
		var hr models.HealthRecord
		if err := rows.Scan(&hr.CheckDate, &hr.Weight, &hr.Length, &hr.Topic, &hr.Notes); err != nil {
			http.Error(w, "error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		records = append(records, hr)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "error iterating rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// return JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(records); err != nil {
		http.Error(w, "failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

// POST snake Health rec
func (s *Server) SnakeHealthPost(w http.ResponseWriter, r *http.Request) {
	//claims
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	userUUID := claims["uuid"].(string)

	//decode json
	defer r.Body.Close()
	req, err := helpers.DecodeBody[models.HealthRecord](w, r)

	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.CheckDate == "" {
		http.Error(w, "check_date is required", http.StatusBadRequest)
		return
	}
	if req.Weight == "" && req.Length == "" && req.Topic == "" && req.Notes == "" {
		http.Error(w, "at least one of weight, length, topic, or notes must be provided", http.StatusBadRequest)
		return
	}

	//get suid
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

	// add health record
	_, err = s.Data_DB.Exec(
		datadb.PostSnakeHealth,
		suid, userUUID, req.Sid, req.CheckDate, req.Weight, req.Length, req.Topic, req.Notes,
	)
	if err != nil {
		http.Error(w, "error inserting health record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"health record added"}`))
}

//UPDATE snake Health rec

func (s *Server) SnakeHealthUpdate(w http.ResponseWriter, r *http.Request) {
	// claims
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	userUUID := claims["uuid"].(string)

	// decode json
	defer r.Body.Close()

	req, err := helpers.DecodeBody[models.HealthRecord](w, r)

	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Sid == "" || req.CheckDate == "" {
		http.Error(w, "sid and check_date are required", http.StatusBadRequest)
		return
	}

	// Collect fields to update
	updates := []string{}
	args := []interface{}{}
	if req.Weight != "" {
		updates = append(updates, "weight = ?")
		args = append(args, req.Weight)
	}
	if req.Length != "" {
		updates = append(updates, "length = ?")
		args = append(args, req.Length)
	}
	if req.Topic != "" {
		updates = append(updates, "topic = ?")
		args = append(args, req.Topic)
	}
	if req.Notes != "" {
		updates = append(updates, "notes = ?")
		args = append(args, req.Notes)
	}

	if len(updates) == 0 {
		http.Error(w, "at least one field (weight, length, topic, or notes) is required", http.StatusBadRequest)
		return
	}

	// get suid
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

	// update record
	query := fmt.Sprintf(
		datadb.UpdateSnakeHealth,
		strings.Join(updates, ", "),
	)
	args = append(args, suid, req.CheckDate)

	res, err := s.Data_DB.Exec(query, args...)
	if err != nil {
		http.Error(w, "error updating health record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "health record not found", http.StatusNotFound)
		return
	}

	// success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"health record updated successfully"}`))

}

// DELETE snake Health rec
func (s *Server) SnakeHealthDelete(w http.ResponseWriter, r *http.Request) {
	//claims
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	userUUID := claims["uuid"].(string)

	//decode json
	req, err := helpers.DecodeBody[models.DeleteSnakeHealth](w, r)

	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Sid == "" || req.CheckDate == "" {
		http.Error(w, "sid and check_date are required", http.StatusBadRequest)
		return
	}

	//get suid

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

	//delete record
	res, err := s.Data_DB.Exec(
		datadb.DeleteSnakeHealth,
		suid, req.CheckDate,
	)
	if err != nil {
		http.Error(w, "error deleting health record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "health record not found", http.StatusNotFound)
		return
	}

	// success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"health record deleted successfully"}`))

}
