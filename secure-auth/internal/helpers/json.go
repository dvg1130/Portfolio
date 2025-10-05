package helpers

import (
	"encoding/json"
	"net/http"
)

func DecodeBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var payload T
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return nil, err
	}
	return &payload, nil

}
