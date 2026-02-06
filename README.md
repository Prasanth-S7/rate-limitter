# Go Rate Limiter

Token bucket rate limiter for Go with per-client (IP/UserID) support.

## Installation

```bash
go get github.com/Prasanth-S7/rate-limitter

```

## Example 

```go
package main

import (
	"net/http"
	"time"

	ratelimiter "github.com/Prasanth-S7/rate-limiter"
	"github.com/Prasanth-S7/rate-limiter/middleware"
)

func main() {
	manager := ratelimiter.NewManager(10, 2, 10*time.Minute)
	limiter := middleware.New(manager, false)

	// Use the limiter in your application, for example, as an HTTP middleware
	http.Handle("/api", limiter.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})))

	http.ListenAndServe(":8080", nil)
}
```

## Manager Configuration

`NewManager(capacity, fillRate, ttl)` creates a rate limiter manager with three parameters:

| Parameter | Type | Description |
|-----------|------|-------------|
| `capacity` | `int` | **Maximum tokens** a client can accumulate. This is the burst capacity - how many requests a client can make in a short spike before being limited. |
| `fillRate` | `float64` | **Tokens added per second**. The steady-state rate limit. A value of `2` means 2 tokens (requests) are added every second. |
| `ttl` | `time.Duration` | **Time-to-live** for inactive client buckets. If a client hasn't made a request in this duration, their bucket is deleted to free memory. |
