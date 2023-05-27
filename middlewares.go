package kafekoding_api

import (
	"context"
	"log"
	"net/http"
	"strings"
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

// authenticationMiddleware is middleware for check authentication for user, by bearer token.
func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			responseJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"status":  "error",
				"message": "Autentikasi dibutuhkan",
			})
			return
		}

		token := strings.Replace(authHeader, "Bearer ", "", -1)

		claims, err := verifyToken(token)
		if err != nil {
			responseJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx := context.WithValue(context.Background(), &userAuth{}, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
