package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"composer/internal/orchestrator"
	"composer/internal/workflow"
)

// buildRunsRouter registers run-related routes
func buildRunsRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/runs", handleGetRuns)
	mux.HandleFunc("GET /api/run/{name}", handleGetRun)
	mux.HandleFunc("POST /api/run/{name}", handlePostRun)
	mux.HandleFunc("GET /api/run/{name}/tasks", handleGetRunTasks)
}

// handleGetRuns returns a list of all runs
func handleGetRuns(w http.ResponseWriter, r *http.Request) {
	runs, err := workflow.ListRuns()
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list runs: %v", err))
		return
	}
	writeData(w, http.StatusOK, runs)
}

// handleGetRun returns a specific run by name
func handleGetRun(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	state, err := workflow.LoadState(name)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("Run not found: %v", err))
		return
	}
	writeData(w, http.StatusOK, state)
}

// handlePostRun creates a new run from a workflow
func handlePostRun(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	var req struct {
		WorkflowName string `json:"workflow_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %v", err))
		return
	}

	if req.WorkflowName == "" {
		writeError(w, http.StatusBadRequest, "workflow_name is required")
		return
	}

	// Load the workflow
	wf, _, err := workflow.LoadWorkflow(req.WorkflowName)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("Workflow not found: %v", err))
		return
	}

	// Create the run
	if err := orchestrator.CreateRun(wf, name); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create run: %v", err))
		return
	}

	// Load and return the created run state
	state, err := workflow.LoadState(name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to load created run: %v", err))
		return
	}

	writeData(w, http.StatusOK, state)
}

// handleGetRunTasks returns all tasks waiting for human intervention
func handleGetRunTasks(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	// Load the run state to get the workflow name
	state, err := workflow.LoadState(name)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("Run not found: %v", err))
		return
	}

	// Load the workflow
	wf, _, err := workflow.LoadWorkflow(state.WorkflowName)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("Workflow not found: %v", err))
		return
	}

	// Get waiting tasks
	tasks, err := orchestrator.ListWaitingTasks(wf, name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list tasks: %v", err))
		return
	}

	writeData(w, http.StatusOK, tasks)
}
