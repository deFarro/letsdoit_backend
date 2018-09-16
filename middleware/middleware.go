package middleware

import (
	"fmt"
	"net/http"
)

// Adapter type for widdlewares
type Adapter func(http.Handler) http.Handler

// Adapt is the function to apple middlewares
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}

	return h
}

// WithLogging is the middleware for logging incoming requests
func WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "OPTIONS" {
			fmt.Printf("Incoming request. Path: %s. Method: %s.\n", r.URL.Path, r.Method)
		}
		h.ServeHTTP(w, r)
	})
}

// WithHeaders is the middleware for writing headers
func WithHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Connection", "keep-alive")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "DNT,Authorization,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")

		h.ServeHTTP(w, r)
	})
}
