package main

import (
	"fmt"
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

func main() {
	tb := NewTokenBucket(10, 2)

	// hit this every 0.1 seconds
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for i := 0; i < 10000; i++ {
		<-ticker.C
		if tb.Allow() {
			fmt.Println("Request", i, "allowed")
		} else {
			fmt.Println("Request", i, "denied")
		}
	}
}
