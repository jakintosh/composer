package ui

import (
	"fmt"
	"net/http"

	"composer/internal/orchestrator"
	"composer/internal/workflow"
)

// BuildRouter creates and configures the UI router.
func (s *Server) BuildRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /", handleDashboard(s.renderer))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(s.static))))
	return mux
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

		tasks, err := orchestrator.ListWaitingTasksByRun(runs)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to load waiting tasks: %v", err), http.StatusInternalServerError)
			return
		}

		data := buildDashboardViewModel(workflows, runs, tasks)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := renderer.Page(w, "pages/dashboard", data); err != nil {
			http.Error(w, fmt.Sprintf("failed to render dashboard: %v", err), http.StatusInternalServerError)
		}
	}
}
