package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Prasanth-S7/rate-limitter/middleware"
	"github.com/Prasanth-S7/rate-limitter/ratelimiter"
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
