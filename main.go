package main

import (
	"log"
	"net/http"
	"time"
)

var b *TokenBucket = NewTokenBucket(Config{
	Tokens:         5,
	Capacity:       5,
	RefillRate:     1,
	RefillInterval: 1 * time.Second,
})

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !b.Allow() {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/api/unprotected", helloWorldHandler)
	http.Handle("/api/protected", rateLimitMiddleware(http.HandlerFunc(helloWorldHandler)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
