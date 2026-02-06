# Go Rate Limiter

Token bucket rate limiter for Go with per-client (IP/UserID) support.

## Installation

```bash
go get github.com/Prasanth-S7/ratelimiter

```

## Example 

```go
package main

import (
    "net/http"
    "time"
    "github.com/Prasanth-S7/ratelimiter"
    "github.com/Prasanth-S7/ratelimiter/middleware"
)

func main() {
    // 10 tokens, 2/sec refill, 10min TTL
    manager := ratelimiter.NewManager(10, 2, 10*time.Minute)
    
    // IP-based limiting
    limiter := middleware.NewHTTP(manager, false)
    
    http.Handle("/", limiter.Handler(yourHandler))
    http.ListenAndServe(":8080", nil)
}
```