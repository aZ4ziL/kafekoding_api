package kafekoding_api

import (
	"encoding/json"
	"net/http"
)

// responseJSON is function to make response json output to client.
func responseJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
