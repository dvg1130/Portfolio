package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/dvg1130/Portfolio/secure-auth/internal/auth"
	"github.com/dvg1130/Portfolio/secure-auth/logs"
	"github.com/google/uuid"
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

			// Safe extraction
			requestID, ok := r.Context().Value(RequestIDKey).(string)
			if !ok || requestID == "" {
				requestID = uuid.NewString()
			}

			var uuid, username, role string
			if claims != nil {
				uuid, _ = claims["uuid"].(string)
				username, _ = claims["username"].(string)
				role, _ = claims["role"].(string)
			}

			// choose level based on status
			if rec.status >= 400 {
				logs.LogEvent(logger, "warn", "request incompleted", r, map[string]interface{}{
					"method":     r.Method,
					"path":       r.URL.Path,
					"status":     rec.status, // recorder’s status code
					"remote_ip":  r.RemoteAddr,
					"latency":    duration, // time.Since(start)
					"uuid":       uuid,
					"user":       username, // extracted from context/JWT
					"role":       role,     // extracted from context/JWT
					"request_id": requestID,
				})
			} else {
				logs.LogEvent(logger, "info", "request completed", r, map[string]interface{}{
					"method":     r.Method,
					"path":       r.URL.Path,
					"status":     rec.status, // recorder’s status code
					"remote_ip":  r.RemoteAddr,
					"latency":    duration, // time.Since(start)
					"uuid":       uuid,
					"user":       username, // extracted from context/JWT
					"role":       role,     // extracted from context/JWT
					"request_id": requestID,
				})
			}
		})
	}
}

func IncomingLoggerMiddleware(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rec, r)

			duration := time.Since(start)
			// pull user/role from context

			// claims := auth.GetClaimsFromContext(r.Context())
			// Safe extraction
			requestID, ok := r.Context().Value(RequestIDKey).(string)
			if !ok || requestID == "" {
				requestID = uuid.NewString()
			}

			// log basic info only
			logFields := map[string]interface{}{
				"method":     r.Method,
				"path":       r.URL.Path,
				"status":     rec.status,
				"remote_ip":  r.RemoteAddr,
				"latency":    duration,
				"request_id": requestID,

				"category": "auth-traffic",
			}

			if rec.status >= 400 {
				logs.LogEvent(logger, "warn", "request completed", r, logFields)
			} else {
				logs.LogEvent(logger, "info", "request completed", r, logFields)
			}
		})
	}
}

type contextKey string

const RequestIDKey contextKey = "requestID"

func RequestIDMiddleware(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.NewString()
			ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

			// optional: log incoming request
			logger.Info("incoming request", map[string]interface{}{
				"category":  "traffic",
				"requestID": requestID,
				"method":    r.Method,
				"path":      r.URL.Path,
			})
			w.Header().Set("X-Request-ID", requestID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
