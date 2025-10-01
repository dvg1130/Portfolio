package middleware

import (
	"database/sql"
	"net/http"

	"github.com/dvg1130/Portfolio/secure-backend/internal/auth"
)

func RecordsLimiter(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// grab claims from context (added by AuthMiddleware)
			claims := auth.GetClaimsFromContext(r.Context())
			if claims == nil {
				http.Error(w, "no claims found", http.StatusUnauthorized)
				return
			}

			uuid, ok1 := claims["uuid"].(string)
			role, ok2 := claims["role"].(string)
			if !ok1 || !ok2 {
				http.Error(w, "invalid claims", http.StatusUnauthorized)
				return
			}

			// count user records
			var count int
			if err := db.QueryRow(
				"SELECT COUNT(*) FROM snakes WHERE owner_uuid = ?",
				uuid,
			).Scan(&count); err != nil {
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}

			// limits by role
			maxByRole := map[string]int{
				"basic": 5,
				"tier1": 15,
				"tier2": 50,
			}
			if limit, ok := maxByRole[role]; ok && count >= limit {
				http.Error(w, "Snake limit reached for your tier", http.StatusForbidden)
				return
			}

			// all good → proceed
			next.ServeHTTP(w, r)
		})
	}
}
