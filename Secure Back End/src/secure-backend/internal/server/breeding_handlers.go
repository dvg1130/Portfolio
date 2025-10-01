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

// GET all SnakeBreeding events
func (s *Server) SnakeBreedGetBySnake(w http.ResponseWriter, r *http.Request) {
	//claims
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	ownerUUID := claims["uuid"].(string)

	//decode json

	defer r.Body.Close()

	req, err := helpers.DecodeBody[models.FindBreedingEvent2](w, r)

	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Sid == "" {
		http.Error(w, "sid and event id required", http.StatusBadRequest)
		return
	}

	//get suid

	var suid string

	err = s.Data_DB.QueryRow(
		datadb.GetSuid,
		ownerUUID, req.Sid,
	).Scan(&suid)
	if err == sql.ErrNoRows {
		http.Error(w, "snake not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error finding snake: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//get all breeding year values
	rows, err := s.Data_DB.Query(
		datadb.GetAllBreedingBySnake,
		ownerUUID, suid, suid, suid, suid, suid,
	)
	if err != nil {
		http.Error(w, "error retrieving breeding records: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	years := []string{}
	for rows.Next() {
		var year string
		if err := rows.Scan(&year); err != nil {
			http.Error(w, "error scanning record: "+err.Error(), http.StatusInternalServerError)
			return
		}
		years = append(years, year)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "row iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//return list json

	resp := models.BreedingEventItem{
		Sid:           req.Sid,
		BreedingYears: years,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func (s *Server) SnakeBreedGetAll(w http.ResponseWriter, r *http.Request) {
	//claims

	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	ownerUUID := claims["uuid"].(string)
	//no json decode
	//
	rows, err := s.Data_DB.Query(datadb.GetAllBreedingsByUser, ownerUUID)
	if err != nil {
		http.Error(w, "error retrieving breeding records: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var breedings []models.BreedingSummary
	for rows.Next() {
		var b models.BreedingSummary
		if err := rows.Scan(&b.BreedingUUID, &b.EventID, &b.Year, &b.Season); err != nil {
			http.Error(w, "error scanning record: "+err.Error(), http.StatusInternalServerError)
			return
		}
		breedings = append(breedings, b)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "row iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breedings)
}

// GET one SnakeBreed event
func (s *Server) SnakeBreedGetOne(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	ownerUUID := claims["uuid"].(string)

	defer r.Body.Close()

	//decode json
	req, err := helpers.DecodeBody[models.FindBreedingEvent](w, r)
	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Sid == "" || req.Event_id == "" {
		http.Error(w, "sid and event_id are required", http.StatusBadRequest)
		return
	}

	//get breeding_uuid
	var breeding_uuid string
	err = s.Data_DB.QueryRow(
		datadb.GetBreedingUuid,
		ownerUUID, req.Event_id,
	).Scan(&breeding_uuid)
	if err == sql.ErrNoRows {
		http.Error(w, "breeding event not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error finding breeding event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//get breeding event data
	var breeding models.BreedingEvent
	err = s.Data_DB.QueryRow(
		datadb.GetBreedingEvent,
		ownerUUID, breeding_uuid,
	).Scan(
		&breeding.Event_id, &breeding.FemaleSid, &breeding.Male1Sid, &breeding.Male2Sid, &breeding.Male3Sid,
		&breeding.Male4Sid, &breeding.BreedingYear, &breeding.BreedingSeason, &breeding.FemaleWeight, &breeding.Male1Weight, &breeding.Male2Weight, &breeding.Male3Weight, &breeding.Male4Weight,
		&breeding.CoolingStart, &breeding.CoolingEnd, &breeding.WarmingStart, &breeding.WarmingEnd, &breeding.PairingDate1, &breeding.PairingDate2, &breeding.PairingDate3, &breeding.PairingDate4,
		&breeding.GravidDate, &breeding.LayDate, &breeding.ClutchSize, &breeding.ClutchSurvive, &breeding.Outcome, &breeding.Notes,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "breeding record not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error retrieving record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//return record
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(breeding); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

// POST SnakeBreed event
func (s *Server) SnakeBreedPost(w http.ResponseWriter, r *http.Request) {
	//claims
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	ownerUUID := claims["uuid"].(string)

	//decode json
	defer r.Body.Close()

	req, err := helpers.DecodeBody[models.BreedingEvent](w, r)

	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.FemaleSid == "" || req.Male1Sid == "" || req.BreedingYear == "" {
		http.Error(w, "sid, mate_sid and breeding_year are required", http.StatusBadRequest)
		return
	}

	//get female suid
	var femaleSuid string
	err = s.Data_DB.QueryRow(
		datadb.GetSuid,
		ownerUUID, req.FemaleSid,
	).Scan(&femaleSuid)
	if err == sql.ErrNoRows {
		http.Error(w, "snake not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error finding snake: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//get male(s) suid
	var male1Suid string
	var male2Suid string
	var male3Suid string
	var male4Suid string

	err = s.Data_DB.QueryRow(
		datadb.GetSuid,
		ownerUUID, req.Male1Sid,
	).Scan(&male1Suid)
	if err == sql.ErrNoRows {
		http.Error(w, "mate snake not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error finding mate snake: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Male2Sid != "" {
		if err := s.Data_DB.QueryRow(datadb.GetSuid, ownerUUID, req.Male2Sid).Scan(&male2Suid); err != nil {
			http.Error(w, "error finding second mate: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if req.Male3Sid != "" {
		if err := s.Data_DB.QueryRow(datadb.GetSuid, ownerUUID, req.Male3Sid).Scan(&male3Suid); err != nil {
			http.Error(w, "error finding third mate: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if req.Male4Sid != "" {
		if err := s.Data_DB.QueryRow(datadb.GetSuid, ownerUUID, req.Male4Sid).Scan(&male4Suid); err != nil {
			http.Error(w, "error finding fourth mate: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	//add breeding event
	_, err = s.Data_DB.Exec(
		datadb.AddBreedingEvent,
		ownerUUID, req.Event_id, femaleSuid, male1Suid, male2Suid, male3Suid, male4Suid,
		req.FemaleSid, req.Male1Sid, req.Male2Sid, req.Male3Sid, req.Male4Sid, req.BreedingYear, req.BreedingSeason, req.FemaleWeight,
		req.Male1Weight, req.Male2Weight, req.Male3Weight, req.Male4Weight,
		req.CoolingStart, req.CoolingEnd, req.WarmingStart, req.WarmingEnd, req.PairingDate1,
		req.PairingDate2, req.PairingDate3, req.PairingDate4, req.GravidDate, req.LayDate, req.ClutchSize,
		req.ClutchSurvive, req.Outcome, req.Notes,
	)
	if err != nil {
		http.Error(w, "error inserting breeding record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"breeding record added"}`))

}

// UPDATE SnakeBreedevent
func (s *Server) SnakeBreedUpdate(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	ownerUUID := claims["uuid"].(string)

	defer r.Body.Close()

	//decode json
	req, err := helpers.DecodeBody[models.BreedingEvent](w, r)
	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Event_id == "" {
		http.Error(w, "breeing event id required", http.StatusBadRequest)
		return
	}

	//get breeding_uuid
	var breeding_uuid string
	err = s.Data_DB.QueryRow(
		datadb.GetBreedingUuid,
		ownerUUID, req.Event_id,
	).Scan(&breeding_uuid)
	if err == sql.ErrNoRows {
		http.Error(w, "breeding event not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error finding breeding event: "+err.Error(), http.StatusInternalServerError)
		return
	}
	//breeding update case
	event := req
	updates := []string{}
	args := []interface{}{}

	if event.FemaleSid != "" {
		updates = append(updates, "female_sid = ?")
		args = append(args, &event.FemaleSid)
	}
	if event.Male1Sid != "" {
		updates = append(updates, "male1_sid = ?")
		args = append(args, &event.Male1Sid)
	}
	if event.Male2Sid != "" {
		updates = append(updates, "male2_sid = ?")
		args = append(args, &event.Male2Sid)
	}
	if event.Male3Sid != "" {
		updates = append(updates, "male3_sid = ?")
		args = append(args, &event.Male3Sid)
	}
	if event.Male4Sid != "" {
		updates = append(updates, "male4_sid = ?")
		args = append(args, &event.Male4Sid)
	}
	if event.BreedingYear != "" {
		updates = append(updates, "breeding_year = ?")
		args = append(args, &event.BreedingYear)
	}
	if event.BreedingSeason != "" {
		updates = append(updates, "breeding_season = ?")
		args = append(args, &event.BreedingSeason)
	}
	if event.FemaleWeight != nil {
		updates = append(updates, "female_weight = ?")
		args = append(args, &event.FemaleWeight)
	}
	if event.Male1Weight != nil {
		updates = append(updates, "male1_weight = ?")
		args = append(args, &event.Male1Weight)
	}
	if event.Male2Weight != nil {
		updates = append(updates, "male2_weight = ?")
		args = append(args, &event.Male2Weight)
	}
	if event.Male3Weight != nil {
		updates = append(updates, "male3_weight = ?")
		args = append(args, &event.Male3Weight)
	}
	if event.Male4Weight != nil {
		updates = append(updates, "male4_weight = ?")
		args = append(args, &event.Male4Weight)
	}
	if event.CoolingStart != nil {
		updates = append(updates, "cooling_start = ?")
		args = append(args, &event.CoolingStart)
	}
	if event.CoolingEnd != nil {
		updates = append(updates, "cooling_end = ?")
		args = append(args, &event.CoolingEnd)
	}
	if event.WarmingStart != nil {
		updates = append(updates, "warming_start = ?")
		args = append(args, &event.WarmingStart)
	}
	if event.WarmingEnd != nil {
		updates = append(updates, "warm_end = ?")
		args = append(args, &event.WarmingEnd)
	}
	if event.PairingDate1 != nil {
		updates = append(updates, "pairing1_date = ?")
		args = append(args, &event.PairingDate1)
	}
	if event.PairingDate2 != nil {
		updates = append(updates, "pairing2_date = ?")
		args = append(args, &event.PairingDate2)
	}
	if event.PairingDate3 != nil {
		updates = append(updates, "pairing3_date = ?")
		args = append(args, &event.PairingDate3)
	}
	if event.PairingDate4 != nil {
		updates = append(updates, "pairing4_date = ?")
		args = append(args, &event.PairingDate4)
	}
	if event.GravidDate != nil {
		updates = append(updates, "gravid_date = ?")
		args = append(args, &event.GravidDate)
	}
	if event.LayDate != nil {
		updates = append(updates, "lay_date = ?")
		args = append(args, &event.LayDate)
	}
	if event.ClutchSize != nil {
		updates = append(updates, "clutch_size = ?")
		args = append(args, &event.ClutchSize)
	}
	if event.ClutchSurvive != nil {
		updates = append(updates, "clutch_survive = ?")
		args = append(args, &event.ClutchSurvive)
	}
	if event.Notes != nil {
		updates = append(updates, "notes = ?")
		args = append(args, &event.Notes)
	}

	if len(updates) == 0 {
		http.Error(w, "no fields to update", http.StatusBadRequest)
		return
	}

	// add WHERE argument
	args = append(args, ownerUUID, breeding_uuid)

	query := fmt.Sprintf(datadb.UpdateBreed,
		strings.Join(updates, ", "))

	if _, err := s.Data_DB.Exec(query, args...); err != nil {
		http.Error(w, "update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"updated"}`))

}

// DELETE SnakeBreed event
func (s *Server) SnakeBreedDelete(w http.ResponseWriter, r *http.Request) {
	//claims
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "no claims found", http.StatusUnauthorized)
		return
	}
	ownerUUID := claims["uuid"].(string)
	//decode json
	defer r.Body.Close()

	req, err := helpers.DecodeBody[models.FindBreedingEvent](w, r)

	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Event_id == "" {
		http.Error(w, "event is required", http.StatusBadRequest)
		return
	}

	//  Find breeding_uuid for this event_id
	var breedingUUID string
	err = s.Data_DB.QueryRow(
		datadb.GetBreedingUuid,
		ownerUUID, req.Event_id,
	).Scan(&breedingUUID)
	if err == sql.ErrNoRows {
		http.Error(w, "breeding record not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error finding breeding record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//delete based on breed_uuid
	res, err := s.Data_DB.Exec(
		datadb.DeleteBreedingEvent,
		ownerUUID, breedingUUID,
	)
	if err != nil {
		http.Error(w, "error deleting snake: "+err.Error(), http.StatusInternalServerError)
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
	w.Write([]byte(`{"message":"breeding record deleted successfully"}`))

}
