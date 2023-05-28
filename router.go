package kafekoding_api

import "net/http"

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir("./media"))))

	// user group
	mux.Handle("/user/register", loggingMiddleware(
		methodMiddleware(http.HandlerFunc(signUpHandler), "POST")))
	mux.Handle("/user/get-token", loggingMiddleware(
		methodMiddleware(http.HandlerFunc(getTokenHandler), "POST")))
	mux.Handle("/user/auth", loggingMiddleware(
		authenticationMiddleware(methodMiddleware(http.HandlerFunc(authHandler), "GET"))))

	// course group
	mux.Handle("/courses", loggingMiddleware(authenticationMiddleware(http.HandlerFunc(courseHandler))))

	return mux
}
