package main

import (
	"fmt"
	"log"
	"net/http"

	"composer/internal/api"
)

func main() {
	addr := "localhost:8080"

	// Build the router from the api package
	mux := api.BuildRouter()

	fmt.Printf("Starting composerd on http://%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
