package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/dvg1130/Portfolio/secure-backend/internal/auth"
	"github.com/dvg1130/Portfolio/secure-backend/logs"
	"go.uber.org/zap"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

type ctxKey string

const loggerKey ctxKey = "logger"

func LoggingMiddleware(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

			// inject logger into context
			ctx := context.WithValue(r.Context(), loggerKey, logger)
			r = r.WithContext(ctx)

			next.ServeHTTP(rec, r)

			duration := time.Since(start)

			// pull user/role from context
			claims := auth.GetClaimsFromContext(r.Context())
			var username, role string
			if claims != nil {
				username, _ = claims["username"].(string)
				role, _ = claims["role"].(string)
			}

			// choose level based on status
			if rec.status >= 400 {
				logs.LogEvent(logger, "warn", "request completed", r, map[string]interface{}{
					"method":    r.Method,
					"path":      r.URL.Path,
					"status":    rec.status, // recorder’s status code
					"remote_ip": r.RemoteAddr,
					"latency":   duration, // time.Since(start)
					"user":      username, // extracted from context/JWT
					"role":      role,     // extracted from context/JWT
				})
			} else {
				logs.LogEvent(logger, "info", "request completed", r, map[string]interface{}{
					"method":    r.Method,
					"path":      r.URL.Path,
					"status":    rec.status, // recorder’s status code
					"remote_ip": r.RemoteAddr,
					"latency":   duration, // time.Since(start)
					"user":      username, // extracted from context/JWT
					"role":      role,     // extracted from context/JWT
				})
			}
		})
	}
}
