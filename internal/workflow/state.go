package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// StepStatus represents the status of a step in a workflow
type StepStatus string

const (
	StatusPending   StepStatus = "pending"
	StatusReady     StepStatus = "ready"
	StatusFailed    StepStatus = "failed"
	StatusSucceeded StepStatus = "succeeded"
)

// StepState represents the state of a single step
type StepState struct {
	Status StepStatus `json:"status"`
}

// RunState represents the complete state of a workflow run
type RunState struct {
	// WorkflowName is the name of the workflow this run belongs to
	WorkflowName string `json:"workflow_name"`
	// StepStates maps step names to their current state
	StepStates map[string]StepState `json:"step_states"`
	// RunName is the name of this run (not persisted to JSON)
	RunName string `json:"-"`
	// artifactPaths maps artifact names to their filesystem paths (not persisted to JSON)
	artifactPaths map[string]string `json:"-"`
}

// NewRunState creates a new run state initialized with pending steps
func NewRunState(workflow *Workflow, runName string) *RunState {
	state := &RunState{
		WorkflowName:  workflow.ID,
		StepStates:    make(map[string]StepState),
		RunName:       runName,
		artifactPaths: make(map[string]string),
	}

	// Initialize all steps as pending
	for _, step := range workflow.Steps {
		state.StepStates[step.Name] = StepState{
			Status: StatusPending,
		}
	}

	return state
}

// Save saves the run state to a JSON file in the run directory
func (rs *RunState) Save() error {
	runDir := GetRunDir(rs.RunName)

	// Create the run directory if it doesn't exist
	if err := os.MkdirAll(runDir, 0755); err != nil {
		return fmt.Errorf("failed to create run directory: %w", err)
	}

	statePath := filepath.Join(runDir, "state.json")
	data, err := json.MarshalIndent(rs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(statePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// LoadState loads the run state from a JSON file in the run directory
func LoadState(runName string) (*RunState, error) {
	runDir := GetRunDir(runName)
	statePath := filepath.Join(runDir, "state.json")

	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state RunState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	// Initialize non-persisted fields
	state.RunName = runName
	state.artifactPaths = make(map[string]string)

	// Scan artifacts directory and populate the map
	artifactsDir := GetArtifactsDir(runName)
	if _, err := os.Stat(artifactsDir); err == nil {
		entries, err := os.ReadDir(artifactsDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read artifacts directory: %w", err)
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				artifactName := entry.Name()
				artifactPath := filepath.Join(artifactsDir, artifactName)
				state.artifactPaths[artifactName] = artifactPath
			}
		}
	}

	return &state, nil
}

// AllStepsCompleted checks if all steps are either succeeded or failed
func (rs *RunState) AllStepsCompleted() bool {
	for _, state := range rs.StepStates {
		if state.Status == StatusPending || state.Status == StatusReady {
			return false
		}
	}
	return true
}

// HasArtifact checks if an artifact with the given name exists
func (rs *RunState) HasArtifact(name string) bool {
	_, exists := rs.artifactPaths[name]
	return exists
}

// ListArtifacts returns a list of all artifact names
func (rs *RunState) ListArtifacts() []string {
	artifacts := make([]string, 0, len(rs.artifactPaths))
	for name := range rs.artifactPaths {
		artifacts = append(artifacts, name)
	}
	return artifacts
}

// ReadArtifact reads the content of a single artifact
func (rs *RunState) ReadArtifact(name string) (string, error) {
	path, exists := rs.artifactPaths[name]
	if !exists {
		return "", fmt.Errorf("artifact %s not found", name)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read artifact %s: %w", name, err)
	}

	return string(data), nil
}

// ReadArtifacts reads multiple artifacts and returns a map of name to content
func (rs *RunState) ReadArtifacts(names []string) (map[string]string, error) {
	artifacts := make(map[string]string)

	for _, name := range names {
		content, err := rs.ReadArtifact(name)
		if err != nil {
			return nil, err
		}
		artifacts[name] = content
	}

	return artifacts, nil
}

// WriteArtifact writes content to an artifact file and updates the artifact registry
func (rs *RunState) WriteArtifact(name, content string) error {
	artifactsDir := GetArtifactsDir(rs.RunName)

	// Create artifacts directory if it doesn't exist
	if err := os.MkdirAll(artifactsDir, 0755); err != nil {
		return fmt.Errorf("failed to create artifacts directory: %w", err)
	}

	artifactPath := filepath.Join(artifactsDir, name)
	if err := os.WriteFile(artifactPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write artifact %s: %w", name, err)
	}

	// Update the artifact registry
	rs.artifactPaths[name] = artifactPath

	return nil
}

// ListRuns returns all runs found in the runs directory
func ListRuns() ([]RunState, error) {
	runsDir := GetRunsDir()

	// Check if runs directory exists
	if _, err := os.Stat(runsDir); os.IsNotExist(err) {
		return []RunState{}, nil
	}

	// Read directory entries
	entries, err := os.ReadDir(runsDir)
	if err != nil {
		return nil, fmt.Errorf("error reading runs directory: %w", err)
	}

	runs := []RunState{}
	for _, entry := range entries {
		// Skip non-directories
		if !entry.IsDir() {
			continue
		}

		runName := entry.Name()

		// Try to load the run state
		state, err := LoadState(runName)
		if err != nil {
			// Skip runs that can't be loaded (might be incomplete or corrupted)
			continue
		}

		runs = append(runs, *state)
	}

	return runs, nil
}
