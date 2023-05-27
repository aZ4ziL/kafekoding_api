package kafekoding_api

import "net/http"

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	// user group
	mux.Handle("/user/register", loggingMiddleware(
		methodMiddleware(http.HandlerFunc(signUpHandler), "POST")))
	mux.Handle("/user/", loggingMiddleware(
		methodMiddleware(http.HandlerFunc(getTokenHandler), "POST")))

	return mux
}
