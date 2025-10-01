package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimiter(rdb *redis.Client, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			ip := r.RemoteAddr
			key := fmt.Sprintf("rate:%s", ip)

			// increment counter
			count, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			if count == 1 {
				// first request in this window → set expiry
				rdb.Expire(ctx, key, window)
			}

			if count > int64(limit) {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
