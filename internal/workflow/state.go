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
}

// NewRunState creates a new run state initialized with pending steps
func NewRunState(workflow *Workflow) *RunState {
	state := &RunState{
		WorkflowName: workflow.ID,
		StepStates:   make(map[string]StepState),
	}

	// Initialize all steps as pending
	for _, step := range workflow.Steps {
		state.StepStates[step.Name] = StepState{
			Status: StatusPending,
		}
	}

	return state
}

// SaveState saves the run state to a JSON file in the run directory
func SaveState(runName string, state *RunState) error {
	runDir := GetRunDir(runName)

	// Create the run directory if it doesn't exist
	if err := os.MkdirAll(runDir, 0755); err != nil {
		return fmt.Errorf("failed to create run directory: %w", err)
	}

	statePath := filepath.Join(runDir, "state.json")
	data, err := json.MarshalIndent(state, "", "  ")
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
