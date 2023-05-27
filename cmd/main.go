package main

import (
	"fmt"
	"net/http"

	"github.com/aZ4ziL/kafekoding_api"
)

func main() {
	fmt.Println("Server is running on :8000...")
	r := kafekoding_api.Router()

	http.ListenAndServe(":8000", r)
}
