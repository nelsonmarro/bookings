package middlewares

import (
	"log"
	"net/http"
)

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request method and URL
		log.Printf("Request Method: %s, Request URL: %s", r.Method, r.URL)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
