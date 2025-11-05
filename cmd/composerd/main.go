package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"composer/internal/api"
	"composer/internal/ui"
)

func main() {
	addr := "0.0.0.0:8080"

	uiServer, err := ui.Init(resolveUIMode())
	if err != nil {
		log.Fatalf("failed to initialize UI: %v", err)
	}

	apiMux := api.BuildRouter()
	uiMux := uiServer.BuildRouter()

	mux := http.NewServeMux()
	mux.Handle("/", uiMux)
	mux.Handle("/api/", apiMux)

	fmt.Printf("Starting composerd on http://%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func resolveUIMode() ui.Mode {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("COMPOSER_ENV")))
	if env == "dev" || env == "development" {
		return ui.ModeDevelopment
	}
	return ui.ModeProduction
}
