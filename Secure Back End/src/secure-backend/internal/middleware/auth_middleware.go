package middleware

import (
	"net/http"
	"strings"

	"github.com/dvg1130/Portfolio/secure-backend/internal/auth"
	"github.com/dvg1130/Portfolio/secure-backend/logs"
	"go.uber.org/zap"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// verify token
		claims, err := auth.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// add claims to request context
		ctx := r.Context()
		ctx = auth.AddClaimsToContext(ctx, claims)
		r = r.WithContext(ctx)

		// continue to the next handler
		next.ServeHTTP(w, r)
	})
}

func RequireRole(role string, logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := auth.GetClaimsFromContext(r.Context())
			if claims == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			username := claims["username"].(string)
			userRole, ok := claims["role"].(string)
			if !ok || userRole != role {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				logs.LogEvent(logger, "warn", "Unauthorized access", r, map[string]interface{}{
					"username": username,
					"role":     role,
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
