package main

import (
	"fmt"
	"net/http"
	"time"

	ratelimiter "github.com/Prasanth-S7/rate-limiter"
	"github.com/Prasanth-S7/rate-limiter/middleware"
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Gophers! You requested: %s\n", r.URL.Path)
}

func main() {
	manager := ratelimiter.NewManager(3, 1, 10*time.Minute)
	limiter := middleware.New(manager, false)

	http.Handle("/", limiter.Handler(http.HandlerFunc(getHandler)))
	http.ListenAndServe(":8000", nil)
}
