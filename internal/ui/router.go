package ui

import (
	"fmt"
	"net/http"

	"composer/internal/orchestrator"
	"composer/internal/workflow"
)

// BuildRouter creates and configures the UI router.
func BuildRouter() (*http.ServeMux, error) {
	renderer, err := NewRenderer(nil)
	if err != nil {
		return nil, fmt.Errorf("prepare renderer: %w", err)
	}

	staticFS, err := staticFileSystem()
	if err != nil {
		return nil, fmt.Errorf("load static assets: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", handleDashboard(renderer))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	return mux, nil
}

func handleDashboard(renderer *Renderer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		waitingTasks, err := orchestrator.ListWaitingTasksByRun(runs)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to load waiting tasks: %v", err), http.StatusInternalServerError)
			return
		}

		data := buildDashboardViewModel(workflows, runs, waitingTasks)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := renderer.Page(w, "pages/dashboard", data); err != nil {
			http.Error(w, fmt.Sprintf("failed to render dashboard: %v", err), http.StatusInternalServerError)
		}
	}
}
