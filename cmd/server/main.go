package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type tokenBucket struct {
	capacity   int
	tokens     int
	fillRate   time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(capacity int, fillRate float64) *tokenBucket {
	return &tokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		fillRate:   time.Duration(float64(time.Second) / fillRate),
		lastRefill: time.Now(),
	}
}

// checking if a token can be consumed
func (tb *tokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// refilling token logic
func (tb *tokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	tokensToAdd := int(elapsed / tb.fillRate)

	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
	}

	tb.lastRefill = tb.lastRefill.Add(time.Duration(tokensToAdd) * tb.fillRate)
}

func (tb *tokenBucket) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !tb.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Gophers! You requested: %s\n", r.URL.Path)
}

func main() {
	tb := NewTokenBucket(10, 2)
	http.Handle("/", tb.middleware(http.HandlerFunc(getHandler)))
}
