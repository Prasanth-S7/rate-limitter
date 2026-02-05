package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Prasanth-S7/rate-limitter/internal/ratelimiter"
	"github.com/Prasanth-S7/rate-limitter/pkg/middleware"
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
