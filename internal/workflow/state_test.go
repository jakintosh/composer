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

	runID := "test-run"
	displayName := "Test Run"
	state := NewRunState(workflow, runID, displayName)

	// Check ID is set
	if state.ID != runID {
		t.Errorf("Expected ID to be %s, got %s", runID, state.ID)
	}
	// Check Name is set
	if state.Name != displayName {
		t.Errorf("Expected Name to be %s, got %s", displayName, state.Name)
	}

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
}

func TestSaveAndLoadState(t *testing.T) {
	// Create a temporary run directory for testing
	tempDir := t.TempDir()
	runID := "test-run"
	displayName := "Test Run"

	// Override GetRunDir for this test
	originalGetwd := os.Getwd
	os.Chdir(tempDir)
	defer os.Chdir(tempDir) // Won't restore but that's ok for tests
	if cwd, err := originalGetwd(); err == nil {
		defer os.Chdir(cwd)
	}

	// Create a state
	state := &RunState{
		ID:            runID,
		Name:          displayName,
		WorkflowName:  "test-workflow",
		artifactPaths: make(map[string]string),
		StepStates: map[string]StepState{
			"step1": {Status: StatusSucceeded},
			"step2": {Status: StatusPending},
			"step3": {Status: StatusFailed},
		},
	}

	// Save the state
	err := state.Save()
	if err != nil {
		t.Fatalf("Failed to save state: %v", err)
	}

	// Verify the file exists
	statePath := filepath.Join(tempDir, ".composer", "runs", runID, "state.json")
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Fatalf("State file was not created at %s", statePath)
	}

	// Load the state back
	loadedState, err := LoadState(runID)
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

	if loadedState.WorkflowName != state.WorkflowName {
		t.Errorf("Loaded workflow ID is %s, expected %s", loadedState.WorkflowName, state.WorkflowName)
	}
	if loadedState.ID != runID {
		t.Errorf("Loaded run ID is %s, expected %s", loadedState.ID, runID)
	}
	if loadedState.Name != displayName {
		t.Errorf("Loaded run name is %s, expected %s", loadedState.Name, displayName)
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

// Artifact tests

func TestWriteAndReadArtifact(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runID := "test-run"
	artifactName := "test-artifact"
	content := "This is test content\nWith multiple lines"

	state := &RunState{
		ID:            runID,
		artifactPaths: make(map[string]string),
	}

	// Write artifact
	err := state.WriteArtifact(artifactName, content)
	if err != nil {
		t.Fatalf("WriteArtifact failed: %v", err)
	}

	// Verify file was created
	artifactPath := filepath.Join(GetArtifactsDir(runID), artifactName)
	if _, err := os.Stat(artifactPath); os.IsNotExist(err) {
		t.Fatalf("Artifact file was not created at %s", artifactPath)
	}

	// Read artifact
	readContent, err := state.ReadArtifact(artifactName)
	if err != nil {
		t.Fatalf("ReadArtifact failed: %v", err)
	}

	// Verify content matches
	if readContent != content {
		t.Errorf("Content mismatch.\nExpected: %s\nGot: %s", content, readContent)
	}
}

func TestHasArtifact(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runID := "test-run"

	state := &RunState{
		ID:            runID,
		artifactPaths: make(map[string]string),
	}

	// Should return false for non-existent artifact
	if state.HasArtifact("nonexistent") {
		t.Error("HasArtifact should return false for non-existent artifact")
	}

	// Write an artifact
	state.WriteArtifact("existing", "content")

	// Should return true for existing artifact
	if !state.HasArtifact("existing") {
		t.Error("HasArtifact should return true for existing artifact")
	}

	// Should still return false for different artifact
	if state.HasArtifact("other") {
		t.Error("HasArtifact should return false for different artifact")
	}
}

func TestListArtifacts(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runID := "test-run"

	state := &RunState{
		ID:            runID,
		artifactPaths: make(map[string]string),
	}

	// Should return empty list when no artifacts exist
	artifacts := state.ListArtifacts()
	if len(artifacts) != 0 {
		t.Errorf("Expected 0 artifacts, got %d", len(artifacts))
	}

	// Write some artifacts
	state.WriteArtifact("artifact1", "content1")
	state.WriteArtifact("artifact2", "content2")
	state.WriteArtifact("artifact3", "content3")

	// List should return all artifacts
	artifacts = state.ListArtifacts()

	if len(artifacts) != 3 {
		t.Errorf("Expected 3 artifacts, got %d", len(artifacts))
	}

	// Verify all artifact names are present
	expectedNames := map[string]bool{
		"artifact1": false,
		"artifact2": false,
		"artifact3": false,
	}

	for _, name := range artifacts {
		if _, exists := expectedNames[name]; exists {
			expectedNames[name] = true
		} else {
			t.Errorf("Unexpected artifact name: %s", name)
		}
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("Expected artifact %s not found in list", name)
		}
	}
}

func TestReadArtifacts(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runID := "test-run"

	state := &RunState{
		ID:            runID,
		artifactPaths: make(map[string]string),
	}

	// Write multiple artifacts
	state.WriteArtifact("doc1", "Content of document 1")
	state.WriteArtifact("doc2", "Content of document 2")
	state.WriteArtifact("doc3", "Content of document 3")

	// Read multiple artifacts
	names := []string{"doc1", "doc3"}
	artifacts, err := state.ReadArtifacts(names)
	if err != nil {
		t.Fatalf("ReadArtifacts failed: %v", err)
	}

	if len(artifacts) != 2 {
		t.Errorf("Expected 2 artifacts, got %d", len(artifacts))
	}

	// Verify contents
	if artifacts["doc1"] != "Content of document 1" {
		t.Errorf("doc1 content mismatch: got %s", artifacts["doc1"])
	}
	if artifacts["doc3"] != "Content of document 3" {
		t.Errorf("doc3 content mismatch: got %s", artifacts["doc3"])
	}

	// doc2 should not be in the map
	if _, exists := artifacts["doc2"]; exists {
		t.Error("doc2 should not be in the returned artifacts")
	}
}

func TestReadArtifactsNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runID := "test-run"

	state := &RunState{
		ID:            runID,
		artifactPaths: make(map[string]string),
	}

	// Try to read a non-existent artifact
	names := []string{"nonexistent"}
	_, err := state.ReadArtifacts(names)
	if err == nil {
		t.Error("ReadArtifacts should return error for non-existent artifact")
	}
}

func TestWriteArtifactCreatesDirectory(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runID := "new-run"

	state := &RunState{
		ID:            runID,
		artifactPaths: make(map[string]string),
	}

	// Verify artifacts directory doesn't exist yet
	artifactsDir := GetArtifactsDir(runID)
	if _, err := os.Stat(artifactsDir); !os.IsNotExist(err) {
		t.Fatalf("Artifacts directory should not exist yet")
	}

	// Write artifact should create the directory
	err := state.WriteArtifact("first", "content")
	if err != nil {
		t.Fatalf("WriteArtifact failed: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(artifactsDir); os.IsNotExist(err) {
		t.Error("WriteArtifact should have created artifacts directory")
	}
}

func TestReadArtifactEmpty(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runID := "test-run"

	state := &RunState{
		ID:            runID,
		artifactPaths: make(map[string]string),
	}

	// Write empty artifact
	err := state.WriteArtifact("empty", "")
	if err != nil {
		t.Fatalf("WriteArtifact failed: %v", err)
	}

	// Read empty artifact
	content, err := state.ReadArtifact("empty")
	if err != nil {
		t.Fatalf("ReadArtifact failed: %v", err)
	}

	if content != "" {
		t.Errorf("Expected empty content, got: %s", content)
	}
}

func TestLoadStatePopulatesArtifactRegistry(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runID := "test-run"

	// Create a state with some artifacts
	state := &RunState{
		WorkflowName:  "test-workflow",
		ID:            runID,
		artifactPaths: make(map[string]string),
		StepStates: map[string]StepState{
			"step1": {Status: StatusSucceeded},
		},
	}

	// Write some artifacts
	state.WriteArtifact("artifact1", "content1")
	state.WriteArtifact("artifact2", "content2")
	state.WriteArtifact("artifact3", "content3")

	// Save the state
	err := state.Save()
	if err != nil {
		t.Fatalf("Failed to save state: %v", err)
	}

	// Load the state back
	loadedState, err := LoadState(runID)
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	// Verify artifact registry was populated
	if !loadedState.HasArtifact("artifact1") {
		t.Error("LoadState should populate artifact registry with artifact1")
	}
	if !loadedState.HasArtifact("artifact2") {
		t.Error("LoadState should populate artifact registry with artifact2")
	}
	if !loadedState.HasArtifact("artifact3") {
		t.Error("LoadState should populate artifact registry with artifact3")
	}

	// Verify we can read artifacts through the loaded state
	content, err := loadedState.ReadArtifact("artifact1")
	if err != nil {
		t.Fatalf("Failed to read artifact from loaded state: %v", err)
	}
	if content != "content1" {
		t.Errorf("Expected 'content1', got '%s'", content)
	}

	// Verify ListArtifacts works
	artifacts := loadedState.ListArtifacts()
	if len(artifacts) != 3 {
		t.Errorf("Expected 3 artifacts in loaded state, got %d", len(artifacts))
	}
}
