package orchestrator

import (
	"os"
	"testing"

	"composer/internal/workflow"
)

func TestCreateRun(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	wf := &workflow.Workflow{
		ID: "test",
		Steps: []workflow.Step{
			{Name: "step1", Output: "out1"},
			{Name: "step2", Inputs: []string{"out1"}, Output: "out2"},
		},
	}

	runName := "test-run"
	err := CreateRun(wf, runName)
	if err != nil {
		t.Fatalf("CreateRun failed: %v", err)
	}

	// Verify state was saved
	state, err := workflow.LoadState(runName)
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	// All steps should be pending
	for _, step := range wf.Steps {
		stepState := state.StepStates[step.Name]
		if stepState.Status != workflow.StatusPending {
			t.Errorf("Step %s should be pending, got %s", step.Name, stepState.Status)
		}
	}

	// No outputs yet
	if len(state.Outputs) != 0 {
		t.Errorf("Expected 0 outputs, got %d", len(state.Outputs))
	}
}

func TestTickWithNoInputs(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	wf := &workflow.Workflow{
		ID: "test",
		Steps: []workflow.Step{
			{Name: "start", Description: "Start step", Output: "started"},
		},
	}

	runName := "test-run"
	CreateRun(wf, runName)

	// First tick should run the step with no inputs
	complete, err := Tick(wf, runName)
	if err != nil {
		t.Fatalf("Tick failed: %v", err)
	}

	if !complete {
		t.Error("Workflow should be complete after first tick")
	}

	// Verify state
	state, _ := workflow.LoadState(runName)
	if state.StepStates["start"].Status != workflow.StatusSucceeded {
		t.Error("Start step should be succeeded")
	}
	if !state.HasOutput("started") {
		t.Error("Should have 'started' output")
	}
}

func TestTickWithDependencies(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	wf := &workflow.Workflow{
		ID: "test",
		Steps: []workflow.Step{
			{Name: "step1", Output: "out1"},
			{Name: "step2", Inputs: []string{"out1"}, Output: "out2"},
			{Name: "step3", Inputs: []string{"out2"}, Output: "out3"},
		},
	}

	runName := "test-run"
	CreateRun(wf, runName)

	// First tick: step1 runs
	complete, err := Tick(wf, runName)
	if err != nil {
		t.Fatalf("First tick failed: %v", err)
	}
	if complete {
		t.Error("Workflow should not be complete after first tick")
	}

	state, _ := workflow.LoadState(runName)
	if state.StepStates["step1"].Status != workflow.StatusSucceeded {
		t.Error("step1 should be succeeded")
	}
	if state.StepStates["step2"].Status != workflow.StatusPending {
		t.Error("step2 should still be pending")
	}

	// Second tick: step2 runs
	complete, err = Tick(wf, runName)
	if err != nil {
		t.Fatalf("Second tick failed: %v", err)
	}
	if complete {
		t.Error("Workflow should not be complete after second tick")
	}

	state, _ = workflow.LoadState(runName)
	if state.StepStates["step2"].Status != workflow.StatusSucceeded {
		t.Error("step2 should be succeeded")
	}
	if state.StepStates["step3"].Status != workflow.StatusPending {
		t.Error("step3 should still be pending")
	}

	// Third tick: step3 runs
	complete, err = Tick(wf, runName)
	if err != nil {
		t.Fatalf("Third tick failed: %v", err)
	}
	if !complete {
		t.Error("Workflow should be complete after third tick")
	}

	state, _ = workflow.LoadState(runName)
	if state.StepStates["step3"].Status != workflow.StatusSucceeded {
		t.Error("step3 should be succeeded")
	}
}

func TestTickWithParallelSteps(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	wf := &workflow.Workflow{
		ID: "test",
		Steps: []workflow.Step{
			{Name: "parallel1", Output: "out1"},
			{Name: "parallel2", Output: "out2"},
			{Name: "parallel3", Output: "out3"},
			{Name: "combine", Inputs: []string{"out1", "out2", "out3"}, Output: "combined"},
		},
	}

	runName := "test-run"
	CreateRun(wf, runName)

	// First tick: all parallel steps run
	complete, err := Tick(wf, runName)
	if err != nil {
		t.Fatalf("First tick failed: %v", err)
	}
	if complete {
		t.Error("Workflow should not be complete after first tick")
	}

	state, _ := workflow.LoadState(runName)
	if state.StepStates["parallel1"].Status != workflow.StatusSucceeded {
		t.Error("parallel1 should be succeeded")
	}
	if state.StepStates["parallel2"].Status != workflow.StatusSucceeded {
		t.Error("parallel2 should be succeeded")
	}
	if state.StepStates["parallel3"].Status != workflow.StatusSucceeded {
		t.Error("parallel3 should be succeeded")
	}

	// Second tick: combine step runs
	complete, err = Tick(wf, runName)
	if err != nil {
		t.Fatalf("Second tick failed: %v", err)
	}
	if !complete {
		t.Error("Workflow should be complete after second tick")
	}

	state, _ = workflow.LoadState(runName)
	if state.StepStates["combine"].Status != workflow.StatusSucceeded {
		t.Error("combine step should be succeeded")
	}
}

func TestTickOnCompleteWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	wf := &workflow.Workflow{
		ID: "test",
		Steps: []workflow.Step{
			{Name: "only", Output: "done"},
		},
	}

	runName := "test-run"
	CreateRun(wf, runName)

	// First tick completes the workflow
	Tick(wf, runName)

	// Second tick on completed workflow should return true
	complete, err := Tick(wf, runName)
	if err != nil {
		t.Fatalf("Tick on complete workflow failed: %v", err)
	}
	if !complete {
		t.Error("Should return complete=true for already completed workflow")
	}
}

func TestFindRunnableSteps(t *testing.T) {
	wf := &workflow.Workflow{
		Steps: []workflow.Step{
			{Name: "step1", Output: "out1"},
			{Name: "step2", Inputs: []string{"out1"}, Output: "out2"},
			{Name: "step3", Inputs: []string{"out1", "out2"}, Output: "out3"},
			{Name: "step4", Inputs: []string{"out3"}, Output: "out4"},
		},
	}

	// Initial state: only step1 should be runnable
	state := workflow.NewRunState(wf)
	runnable := findRunnableSteps(wf, state)
	if len(runnable) != 1 || runnable[0].Name != "step1" {
		t.Errorf("Expected only step1 to be runnable, got %v", runnable)
	}

	// After step1 completes
	state.StepStates["step1"] = workflow.StepState{Status: workflow.StatusSucceeded}
	state.AddOutput("out1")
	runnable = findRunnableSteps(wf, state)
	if len(runnable) != 1 || runnable[0].Name != "step2" {
		t.Errorf("Expected only step2 to be runnable, got %v", runnable)
	}

	// After step2 completes
	state.StepStates["step2"] = workflow.StepState{Status: workflow.StatusSucceeded}
	state.AddOutput("out2")
	runnable = findRunnableSteps(wf, state)
	if len(runnable) != 1 || runnable[0].Name != "step3" {
		t.Errorf("Expected only step3 to be runnable, got %v", runnable)
	}

	// After step3 completes
	state.StepStates["step3"] = workflow.StepState{Status: workflow.StatusSucceeded}
	state.AddOutput("out3")
	runnable = findRunnableSteps(wf, state)
	if len(runnable) != 1 || runnable[0].Name != "step4" {
		t.Errorf("Expected only step4 to be runnable, got %v", runnable)
	}

	// After all complete
	state.StepStates["step4"] = workflow.StepState{Status: workflow.StatusSucceeded}
	runnable = findRunnableSteps(wf, state)
	if len(runnable) != 0 {
		t.Errorf("Expected no runnable steps, got %v", runnable)
	}
}

func TestFindRunnableStepsWithMultipleInputs(t *testing.T) {
	wf := &workflow.Workflow{
		Steps: []workflow.Step{
			{Name: "a", Output: "out_a"},
			{Name: "b", Output: "out_b"},
			{Name: "c", Inputs: []string{"out_a", "out_b"}, Output: "out_c"},
		},
	}

	state := workflow.NewRunState(wf)

	// Both a and b should be runnable
	runnable := findRunnableSteps(wf, state)
	if len(runnable) != 2 {
		t.Errorf("Expected 2 runnable steps, got %d", len(runnable))
	}

	// After only a completes, c should not be runnable
	state.StepStates["a"] = workflow.StepState{Status: workflow.StatusSucceeded}
	state.AddOutput("out_a")
	runnable = findRunnableSteps(wf, state)
	if len(runnable) != 1 || runnable[0].Name != "b" {
		t.Errorf("Expected only b to be runnable, got %v", runnable)
	}

	// After b completes, c should be runnable
	state.StepStates["b"] = workflow.StepState{Status: workflow.StatusSucceeded}
	state.AddOutput("out_b")
	runnable = findRunnableSteps(wf, state)
	if len(runnable) != 1 || runnable[0].Name != "c" {
		t.Errorf("Expected only c to be runnable, got %v", runnable)
	}
}
