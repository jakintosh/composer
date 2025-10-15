package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"composer/internal/workflow"
)

// buildWorkflowsRouter registers workflow-related routes
func buildWorkflowsRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/workflows", handleGetWorkflows)
	mux.HandleFunc("GET /api/workflow/{id}", handleGetWorkflow)
	mux.HandleFunc("POST /api/workflow/{id}", handlePostWorkflow)
}

// handleGetWorkflows returns a list of all workflows
func handleGetWorkflows(w http.ResponseWriter, r *http.Request) {
	workflows, err := workflow.ListWorkflows()
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list workflows: %v", err))
		return
	}
	writeData(w, http.StatusOK, workflows)
}

// handleGetWorkflow returns a specific workflow by ID
func handleGetWorkflow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	wf, _, err := workflow.LoadWorkflow(id)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("Workflow not found: %v", err))
		return
	}
	writeData(w, http.StatusOK, wf)
}

// handlePostWorkflow creates or updates a workflow
func handlePostWorkflow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var wf workflow.Workflow
	if err := json.NewDecoder(r.Body).Decode(&wf); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %v", err))
		return
	}

	// Set the ID from the URL
	wf.ID = id

	if err := workflow.SaveWorkflow(&wf); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save workflow: %v", err))
		return
	}

	writeData(w, http.StatusOK, wf)
}
