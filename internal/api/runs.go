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
	mux.HandleFunc("GET /api/runs/tasks", handleGetRunsTasks)
	mux.HandleFunc("GET /api/run/{id}", handleGetRun)
	mux.HandleFunc("POST /api/run/{id}", handlePostRun)
	mux.HandleFunc("GET /api/run/{id}/tasks", handleGetRunTasks)
	mux.HandleFunc("POST /api/run/{id}/tick", handlePostRunTick)
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
	id := r.PathValue("id")

	state, err := workflow.LoadState(id)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("Run not found: %v", err))
		return
	}
	writeData(w, http.StatusOK, state)
}

// handlePostRun creates a new run from a workflow
func handlePostRun(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req struct {
		WorkflowId     string `json:"workflow_id"`
		RunDisplayName string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %v", err))
		return
	}

	if req.WorkflowId == "" {
		writeError(w, http.StatusBadRequest, "workflow_name is required")
		return
	}
	if req.RunDisplayName == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	// Load the workflow
	wf, _, err := workflow.LoadWorkflow(req.WorkflowId)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("Workflow not found: %v", err))
		return
	}

	// Create the run
	if err := orchestrator.CreateRun(wf, id, req.RunDisplayName); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create run: %v", err))
		return
	}

	// Load and return the created run state
	state, err := workflow.LoadState(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to load created run: %v", err))
		return
	}

	writeData(w, http.StatusOK, state)
}

// handleGetRunTasks returns all tasks waiting for human intervention
func handleGetRunTasks(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// Load the run state to get the workflow ID
	state, err := workflow.LoadState(id)
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
	tasks, err := orchestrator.ListWaitingTasks(wf, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list tasks: %v", err))
		return
	}

	writeData(w, http.StatusOK, tasks)
}

// handleGetRunsTasks returns waiting tasks grouped by run name
func handleGetRunsTasks(w http.ResponseWriter, r *http.Request) {
	runs, err := workflow.ListRuns()
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list runs: %v", err))
		return
	}

	result, err := orchestrator.ListWaitingTasksByRun(runs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list waiting tasks: %v", err))
		return
	}

	writeData(w, http.StatusOK, result)
}

// handlePostRunTick executes a single tick for the specified run
func handlePostRunTick(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// Load current run state to identify its workflow
	state, err := workflow.LoadState(id)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("Run not found: %v", err))
		return
	}

	// Load workflow associated with this run
	wf, _, err := workflow.LoadWorkflow(state.WorkflowName)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("Workflow not found: %v", err))
		return
	}

	// Execute tick
	complete, err := orchestrator.Tick(wf, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to tick run: %v", err))
		return
	}

	// Reload state to include updates from tick
	updatedState, err := workflow.LoadState(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to load updated run state: %v", err))
		return
	}

	writeData(w, http.StatusOK, struct {
		Complete bool               `json:"complete"`
		State    *workflow.RunState `json:"state"`
	}{
		Complete: complete,
		State:    updatedState,
	})
}
