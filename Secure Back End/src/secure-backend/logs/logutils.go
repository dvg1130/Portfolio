package logs

import (
	"net/http"
	"time"

	"github.com/dvg1130/Portfolio/secure-backend/internal/helpers"
	"go.uber.org/zap"
)

func LogEvent(logger *zap.SugaredLogger, level string, msg string, r *http.Request, fields map[string]interface{}) {
	base := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"ip":        helpers.ClientIP(r),
		"path":      r.URL.Path,
	}
	for k, v := range fields {
		base[k] = v
	}

	args := make([]interface{}, 0, len(base)*2)
	for k, v := range base {
		args = append(args, k, v)
	}

	switch level {
	case "warn":
		logger.Warnw(msg, args...)
	case "error":
		logger.Errorw(msg, args...)
	case "info":
		logger.Infow(msg, args...)
	default:
		logger.Debugw(msg, args...)
	}
}
