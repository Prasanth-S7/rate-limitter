package middleware

import (
	"net/http"

	"github.com/Prasanth-S7/rate-limitter/internal/ratelimiter"
)

type RateLimiter struct {
	limiter *ratelimiter.TokenBucket
}

type Middleware struct {
	manager   *ratelimiter.Manager
	useUserID bool
}

func New(manager *ratelimiter.Manager, useUserID bool) *Middleware {
	return &Middleware{
		manager:   manager,
		useUserID: useUserID,
	}
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientKey := m.manager.ExtractClientKey(r, m.useUserID)
		bucket := m.manager.GetBucket(clientKey)

		if !bucket.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
