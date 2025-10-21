package ui

import (
	"fmt"
	"net/http"

	"composer/internal/workflow"
)

// BuildRouter creates and configures the UI router.
func BuildRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleDashboard)
	return mux
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	workflows, err := workflow.ListWorkflows()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load workflows: %v", err), http.StatusInternalServerError)
		return
	}

	runs, err := workflow.ListRuns()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load runs: %v", err), http.StatusInternalServerError)
		return
	}

	data := buildDashboardData(workflows, runs)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := renderDashboard(w, data); err != nil {
		http.Error(w, fmt.Sprintf("failed to render dashboard: %v", err), http.StatusInternalServerError)
		return
	}
}
