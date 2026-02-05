package middleware

import (
	"net/http"

	"github.com/Prasanth-S7/rate-limitter/internal/ratelimiter"
)

type RateLimiter struct {
	limiter *ratelimiter.TokenBucket
}

func New(tb *ratelimiter.TokenBucket) *RateLimiter {
	return &RateLimiter{limiter: tb}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
