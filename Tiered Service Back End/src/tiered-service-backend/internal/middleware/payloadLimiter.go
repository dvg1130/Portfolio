package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tiered-service-backend/internal/auth"
)

func PayloadLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaimsFromContext(r.Context())
		if claims == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		role, _ := claims["role"].(string)

		// role limits
		var maxAllowed int
		switch role {
		case "user":
			maxAllowed = 10
		case "tier1":
			maxAllowed = 100
		case "tier2":
			maxAllowed = 500
		case "admin":
			maxAllowed = -1 // unlimited
		default:
			http.Error(w, "Unknown role", http.StatusForbidden)
			return
		}

		// parse JSON body
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// extract "amount"
		amount, ok := body["amount"].(float64)
		if !ok {
			http.Error(w, "Missing or invalid 'amount'", http.StatusBadRequest)
			return
		}

		if maxAllowed != -1 && int(amount) > maxAllowed {
			http.Error(w, fmt.Sprintf("Amount exceeds limit for role %s (max %d)", role, maxAllowed), http.StatusForbidden)
			return
		}

		// re-encode the body so the next handler can still use it
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			http.Error(w, "Failed to re-encode body", http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(buf)

		next.ServeHTTP(w, r)
	})
}
