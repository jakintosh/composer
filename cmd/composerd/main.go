package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"composer/internal/orchestrator"
	"composer/internal/workflow"
)

func main() {
	addr := "localhost:8080"

	mux := http.NewServeMux()

	// Root handler
	mux.HandleFunc("GET /", handleRoot)

	// Workflow routes
	mux.HandleFunc("GET /api/workflows", handleGetWorkflows)
	mux.HandleFunc("GET /api/workflow/{id}", handleGetWorkflow)
	mux.HandleFunc("POST /api/workflow/{id}", handlePostWorkflow)

	// Run routes
	mux.HandleFunc("GET /api/runs", handleGetRuns)
	mux.HandleFunc("GET /api/run/{name}", handleGetRun)
	mux.HandleFunc("POST /api/run/{name}", handlePostRun)
	mux.HandleFunc("GET /api/run/{name}/tasks", handleGetRunTasks)

	fmt.Printf("Starting composerd on http://%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// GET /
func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<html><body><h1>Coming soon</h1></body></html>")
}

// GET /api/workflows
func handleGetWorkflows(w http.ResponseWriter, r *http.Request) {
	workflows, err := workflow.ListWorkflows()
	if err != nil {
		writeError(w, fmt.Sprintf("Failed to list workflows: %v", err), http.StatusInternalServerError)
		return
	}
	writeJSON(w, workflows)
}

// GET /api/workflow/{id}
func handleGetWorkflow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	wf, _, err := workflow.LoadWorkflow(id)
	if err != nil {
		writeError(w, fmt.Sprintf("Workflow not found: %v", err), http.StatusNotFound)
		return
	}
	writeJSON(w, wf)
}

// POST /api/workflow/{id}
func handlePostWorkflow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var wf workflow.Workflow
	if err := json.NewDecoder(r.Body).Decode(&wf); err != nil {
		writeError(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Set the ID from the URL
	wf.ID = id

	if err := workflow.SaveWorkflow(&wf); err != nil {
		writeError(w, fmt.Sprintf("Failed to save workflow: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, wf)
}

// GET /api/runs
func handleGetRuns(w http.ResponseWriter, r *http.Request) {
	runs, err := workflow.ListRuns()
	if err != nil {
		writeError(w, fmt.Sprintf("Failed to list runs: %v", err), http.StatusInternalServerError)
		return
	}
	writeJSON(w, runs)
}

// GET /api/run/{name}
func handleGetRun(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	state, err := workflow.LoadState(name)
	if err != nil {
		writeError(w, fmt.Sprintf("Run not found: %v", err), http.StatusNotFound)
		return
	}
	writeJSON(w, state)
}

// POST /api/run/{name}
func handlePostRun(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	var req struct {
		WorkflowName string `json:"workflow_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if req.WorkflowName == "" {
		writeError(w, "workflow_name is required", http.StatusBadRequest)
		return
	}

	// Load the workflow
	wf, _, err := workflow.LoadWorkflow(req.WorkflowName)
	if err != nil {
		writeError(w, fmt.Sprintf("Workflow not found: %v", err), http.StatusNotFound)
		return
	}

	// Create the run
	if err := orchestrator.CreateRun(wf, name); err != nil {
		writeError(w, fmt.Sprintf("Failed to create run: %v", err), http.StatusInternalServerError)
		return
	}

	// Load and return the created run state
	state, err := workflow.LoadState(name)
	if err != nil {
		writeError(w, fmt.Sprintf("Failed to load created run: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, state)
}

// GET /api/run/{name}/tasks
func handleGetRunTasks(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	// Load the run state to get the workflow name
	state, err := workflow.LoadState(name)
	if err != nil {
		writeError(w, fmt.Sprintf("Run not found: %v", err), http.StatusNotFound)
		return
	}

	// Load the workflow
	wf, _, err := workflow.LoadWorkflow(state.WorkflowName)
	if err != nil {
		writeError(w, fmt.Sprintf("Workflow not found: %v", err), http.StatusNotFound)
		return
	}

	// Get waiting tasks
	tasks, err := orchestrator.ListWaitingTasks(wf, name)
	if err != nil {
		writeError(w, fmt.Sprintf("Failed to list tasks: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, tasks)
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// writeError writes a JSON error response
func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
