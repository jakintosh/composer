package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	addr := "localhost:8080"

	http.HandleFunc("/", handleRoot)

	fmt.Printf("Starting composerd on http://%s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "hello world!",
	})
}
