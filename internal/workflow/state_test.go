package workflow

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRunState(t *testing.T) {
	workflow := &Workflow{
		ID:          "test",
		Description: "test workflow",
		Steps: []Step{
			{Name: "step1", Output: "out1"},
			{Name: "step2", Inputs: []string{"out1"}, Output: "out2"},
			{Name: "step3", Inputs: []string{"out2"}, Output: "out3"},
		},
	}

	state := NewRunState(workflow)

	// Check all steps are initialized as pending
	if len(state.StepStates) != 3 {
		t.Errorf("Expected 3 step states, got %d", len(state.StepStates))
	}

	for _, step := range workflow.Steps {
		stepState, exists := state.StepStates[step.Name]
		if !exists {
			t.Errorf("Step %s not found in state", step.Name)
		}
		if stepState.Status != StatusPending {
			t.Errorf("Step %s status should be pending, got %s", step.Name, stepState.Status)
		}
	}

	// Check outputs is empty
	if len(state.Outputs) != 0 {
		t.Errorf("Expected 0 outputs, got %d", len(state.Outputs))
	}
}

func TestSaveAndLoadState(t *testing.T) {
	// Create a temporary run directory for testing
	tempDir := t.TempDir()
	runName := "test-run"

	// Override GetRunDir for this test
	originalGetwd := os.Getwd
	os.Chdir(tempDir)
	defer os.Chdir(tempDir) // Won't restore but that's ok for tests
	if cwd, err := originalGetwd(); err == nil {
		defer os.Chdir(cwd)
	}

	// Create a state
	state := &RunState{
		StepStates: map[string]StepState{
			"step1": {Status: StatusSucceeded},
			"step2": {Status: StatusPending},
			"step3": {Status: StatusFailed},
		},
		Outputs: []string{"output1", "output2"},
	}

	// Save the state
	err := SaveState(runName, state)
	if err != nil {
		t.Fatalf("Failed to save state: %v", err)
	}

	// Verify the file exists
	statePath := filepath.Join(tempDir, ".composer", "runs", runName, "state.json")
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Fatalf("State file was not created at %s", statePath)
	}

	// Load the state back
	loadedState, err := LoadState(runName)
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	// Verify the loaded state matches
	if len(loadedState.StepStates) != len(state.StepStates) {
		t.Errorf("Loaded state has %d steps, expected %d", len(loadedState.StepStates), len(state.StepStates))
	}

	for name, stepState := range state.StepStates {
		loadedStepState, exists := loadedState.StepStates[name]
		if !exists {
			t.Errorf("Step %s not found in loaded state", name)
		}
		if loadedStepState.Status != stepState.Status {
			t.Errorf("Step %s status is %s, expected %s", name, loadedStepState.Status, stepState.Status)
		}
	}

	if len(loadedState.Outputs) != len(state.Outputs) {
		t.Errorf("Loaded state has %d outputs, expected %d", len(loadedState.Outputs), len(state.Outputs))
	}

	for i, output := range state.Outputs {
		if loadedState.Outputs[i] != output {
			t.Errorf("Output[%d] is %s, expected %s", i, loadedState.Outputs[i], output)
		}
	}
}

func TestHasOutput(t *testing.T) {
	state := &RunState{
		Outputs: []string{"output1", "output2", "output3"},
	}

	tests := []struct {
		output   string
		expected bool
	}{
		{"output1", true},
		{"output2", true},
		{"output3", true},
		{"output4", false},
		{"", false},
	}

	for _, tt := range tests {
		result := state.HasOutput(tt.output)
		if result != tt.expected {
			t.Errorf("HasOutput(%s) = %v, expected %v", tt.output, result, tt.expected)
		}
	}
}

func TestAddOutput(t *testing.T) {
	state := &RunState{
		Outputs: []string{"output1"},
	}

	// Add a new output
	state.AddOutput("output2")
	if len(state.Outputs) != 2 {
		t.Errorf("Expected 2 outputs, got %d", len(state.Outputs))
	}
	if !state.HasOutput("output2") {
		t.Error("output2 should be in outputs")
	}

	// Try to add a duplicate
	state.AddOutput("output1")
	if len(state.Outputs) != 2 {
		t.Errorf("Expected 2 outputs after duplicate add, got %d", len(state.Outputs))
	}

	// Add another new output
	state.AddOutput("output3")
	if len(state.Outputs) != 3 {
		t.Errorf("Expected 3 outputs, got %d", len(state.Outputs))
	}
}

func TestAllStepsCompleted(t *testing.T) {
	tests := []struct {
		name       string
		stepStates map[string]StepState
		expected   bool
	}{
		{
			name: "all succeeded",
			stepStates: map[string]StepState{
				"step1": {Status: StatusSucceeded},
				"step2": {Status: StatusSucceeded},
			},
			expected: true,
		},
		{
			name: "all failed",
			stepStates: map[string]StepState{
				"step1": {Status: StatusFailed},
				"step2": {Status: StatusFailed},
			},
			expected: true,
		},
		{
			name: "mixed succeeded and failed",
			stepStates: map[string]StepState{
				"step1": {Status: StatusSucceeded},
				"step2": {Status: StatusFailed},
			},
			expected: true,
		},
		{
			name: "one pending",
			stepStates: map[string]StepState{
				"step1": {Status: StatusSucceeded},
				"step2": {Status: StatusPending},
			},
			expected: false,
		},
		{
			name: "all pending",
			stepStates: map[string]StepState{
				"step1": {Status: StatusPending},
				"step2": {Status: StatusPending},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &RunState{StepStates: tt.stepStates}
			result := state.AllStepsCompleted()
			if result != tt.expected {
				t.Errorf("AllStepsCompleted() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
