package middleware

import (
	"context"
	"net/http"
	"tiered-service-backend/internal/auth"
	"time"

	"go.uber.org/zap"
)

// statusRecorder wraps ResponseWriter to capture status codes
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

func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	sugar := logger.Sugar()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

			// inject logger into context
			ctx := context.WithValue(r.Context(), loggerKey, sugar)
			r = r.WithContext(ctx)

			next.ServeHTTP(rec, r)

			duration := time.Since(start)

			//  pull user/role from context
			claims := auth.GetClaimsFromContext(r.Context())

			var user, role string
			if claims != nil {
				user, _ = claims["username"].(string)
				role, _ = claims["role"].(string)
			}

			sugar.Infow("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", rec.status),
				zap.String("remote_ip", r.RemoteAddr),
				zap.Duration("latency", duration),
				zap.String("user", user),
				zap.String("role", role),
			)
		})
	}
}
