package kafekoding_api

import (
	"log"
	"net/http"
)

// loggingMiddleware is function to print a output for every request.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("URL %s - %s\n", r.URL.String(), r.Method)
		next.ServeHTTP(w, r)
	})
}

// methodMiddleware is middleware to handle method type for request.
func methodMiddleware(next http.Handler, method string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			next.ServeHTTP(w, r)
			return
		} else {
			responseJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
				"status":  "method_not_allowed",
				"message": "Method yang anda gunakan tidak ditemukan.",
			})
			return
		}
	})
}
