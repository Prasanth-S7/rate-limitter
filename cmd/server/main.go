package main

import (
	"fmt"
	"net/http"

	"github.com/Prasanth-S7/rate-limitter/internal/ratelimiter"
	"github.com/Prasanth-S7/rate-limitter/pkg/middleware"
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Gophers! You requested: %s\n", r.URL.Path)
}

func main() {
	tb := ratelimiter.NewTokenBucket(10, 2)
	rl := middleware.New(tb)

	http.Handle("/", rl.Middleware(http.HandlerFunc(getHandler)))
	http.ListenAndServe(":8000", nil)
}
