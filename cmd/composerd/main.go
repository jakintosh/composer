package main

import (
	"fmt"
	"log"
	"net/http"

	"composer/internal/api"
	"composer/internal/ui"
)

func main() {
	addr := "localhost:8080"

	apiMux := api.BuildRouter()
	uiMux := ui.BuildRouter()

	mux := http.NewServeMux()
	mux.Handle("/", uiMux)
	mux.Handle("/api/", apiMux)

	fmt.Printf("Starting composerd on http://%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
