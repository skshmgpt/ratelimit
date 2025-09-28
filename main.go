package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var l *RateLimiter = NewRateLimiter(Config{
	capacity:       5,
	refillRate:     1,
	refillInterval: 1 * time.Second,
	redisAddr:      "localhost:6379",
})

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := r.Header.Get("X-Api-Key")
		if id == "" {
			http.Error(w, "api key expected", http.StatusBadRequest)
			return
		}

		if !l.Allow(id) {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/api/unprotected", helloWorldHandler)
	http.Handle("/api/protected", rateLimitMiddleware(http.HandlerFunc(helloWorldHandler)))
	http.HandleFunc("/api/quota", func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Api-Key")
		if id == "" {
			http.Error(w, "api key expected", http.StatusBadRequest)
			return
		}
		tokens := l.GetQuota(id)
		if tokens == "" {
			http.Error(w, "error getting quota", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"tokens": tokens})
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
