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