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
			{Name: "step1", Content: "initial", Output: "out1"},
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

	// No artifacts yet
	artifacts := state.ListArtifacts()
	if len(artifacts) != 0 {
		t.Errorf("Expected 0 artifacts, got %d", len(artifacts))
	}
}

func TestTickWithNoInputs(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	wf := &workflow.Workflow{
		ID: "test",
		Steps: []workflow.Step{
			{Name: "start", Description: "Start step", Content: "initial content", Output: "started"},
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

	// Reload state to access artifacts
	state, _ = workflow.LoadState(runName)

	// Verify artifact was created
	if !state.HasArtifact("started") {
		t.Error("Should have 'started' artifact")
	}

	// Verify artifact content
	content, err := state.ReadArtifact("started")
	if err != nil {
		t.Fatalf("Failed to read artifact: %v", err)
	}
	if content != "initial content" {
		t.Errorf("Artifact content should be 'initial content', got '%s'", content)
	}
}

func TestTickWithDependencies(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	wf := &workflow.Workflow{
		ID: "test",
		Steps: []workflow.Step{
			{Name: "step1", Content: "step1 content", Output: "out1"},
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

	// Reload state to access artifacts
	state, _ = workflow.LoadState(runName)

	// Verify artifact content is from step1
	content, _ := state.ReadArtifact("out1")
	if content != "step1 content" {
		t.Errorf("out1 should contain 'step1 content', got '%s'", content)
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

	// Reload state to access artifacts
	state, _ = workflow.LoadState(runName)

	// Verify artifact content is concatenated from step1
	content, _ = state.ReadArtifact("out2")
	if content != "step1 content" {
		t.Errorf("out2 should contain 'step1 content', got '%s'", content)
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

	// Reload state to access artifacts
	state, _ = workflow.LoadState(runName)

	// Verify final artifact content
	content, _ = state.ReadArtifact("out3")
	if content != "step1 content" {
		t.Errorf("out3 should contain 'step1 content', got '%s'", content)
	}
}

func TestTickWithParallelSteps(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	wf := &workflow.Workflow{
		ID: "test",
		Steps: []workflow.Step{
			{Name: "parallel1", Content: "data1", Output: "out1"},
			{Name: "parallel2", Content: "data2", Output: "out2"},
			{Name: "parallel3", Content: "data3", Output: "out3"},
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
			{Name: "only", Content: "done content", Output: "done"},
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
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runName := "test-run"
	wf := &workflow.Workflow{
		Steps: []workflow.Step{
			{Name: "step1", Content: "initial", Output: "out1"},
			{Name: "step2", Inputs: []string{"out1"}, Output: "out2"},
			{Name: "step3", Inputs: []string{"out1", "out2"}, Output: "out3"},
			{Name: "step4", Inputs: []string{"out3"}, Output: "out4"},
		},
	}

	// Initial state: only step1 should be runnable
	state := workflow.NewRunState(wf, runName)
	runnable := findRunnableSteps(wf, state)
	if len(runnable) != 1 || runnable[0].Name != "step1" {
		t.Errorf("Expected only step1 to be runnable, got %v", runnable)
	}

	// After step1 completes
	state.StepStates["step1"] = workflow.StepState{Status: workflow.StatusSucceeded}
	state.WriteArtifact("out1", "artifact1 content")
	runnable = findRunnableSteps(wf, state)
	if len(runnable) != 1 || runnable[0].Name != "step2" {
		t.Errorf("Expected only step2 to be runnable, got %v", runnable)
	}

	// After step2 completes
	state.StepStates["step2"] = workflow.StepState{Status: workflow.StatusSucceeded}
	state.WriteArtifact("out2", "artifact2 content")
	runnable = findRunnableSteps(wf, state)
	if len(runnable) != 1 || runnable[0].Name != "step3" {
		t.Errorf("Expected only step3 to be runnable, got %v", runnable)
	}

	// After step3 completes
	state.StepStates["step3"] = workflow.StepState{Status: workflow.StatusSucceeded}
	state.WriteArtifact("out3", "artifact3 content")
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
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runName := "test-run"
	wf := &workflow.Workflow{
		Steps: []workflow.Step{
			{Name: "a", Content: "content a", Output: "out_a"},
			{Name: "b", Content: "content b", Output: "out_b"},
			{Name: "c", Inputs: []string{"out_a", "out_b"}, Output: "out_c"},
		},
	}

	state := workflow.NewRunState(wf, runName)

	// Both a and b should be runnable
	runnable := findRunnableSteps(wf, state)
	if len(runnable) != 2 {
		t.Errorf("Expected 2 runnable steps, got %d", len(runnable))
	}

	// After only a completes, c should not be runnable
	state.StepStates["a"] = workflow.StepState{Status: workflow.StatusSucceeded}
	state.WriteArtifact("out_a", "artifact a")
	runnable = findRunnableSteps(wf, state)
	if len(runnable) != 1 || runnable[0].Name != "b" {
		t.Errorf("Expected only b to be runnable, got %v", runnable)
	}

	// After b completes, c should be runnable
	state.StepStates["b"] = workflow.StepState{Status: workflow.StatusSucceeded}
	state.WriteArtifact("out_b", "artifact b")
	runnable = findRunnableSteps(wf, state)
	if len(runnable) != 1 || runnable[0].Name != "c" {
		t.Errorf("Expected only c to be runnable, got %v", runnable)
	}
}

func TestHumanHandlerStep(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	wf := &workflow.Workflow{
		ID: "test",
		Steps: []workflow.Step{
			{Name: "automated", Handler: "tool", Content: "auto content", Output: "auto-out"},
			{Name: "manual", Handler: "human", Prompt: "Please review the data", Inputs: []string{"auto-out"}, Output: "manual-out"},
		},
	}

	runName := "test-run"
	CreateRun(wf, runName)

	// First tick: automated step runs
	complete, err := Tick(wf, runName)
	if err != nil {
		t.Fatalf("First tick failed: %v", err)
	}
	if complete {
		t.Error("Workflow should not be complete after first tick")
	}

	state, _ := workflow.LoadState(runName)
	if state.StepStates["automated"].Status != workflow.StatusSucceeded {
		t.Error("Automated step should be succeeded")
	}

	// Second tick: human step should transition to ready, not succeed
	complete, err = Tick(wf, runName)
	if err != nil {
		t.Fatalf("Second tick failed: %v", err)
	}
	if complete {
		t.Error("Workflow should not be complete with ready task")
	}

	state, _ = workflow.LoadState(runName)
	if state.StepStates["manual"].Status != workflow.StatusReady {
		t.Errorf("Manual step should be ready, got %s", state.StepStates["manual"].Status)
	}

	// Manual step artifact should not exist yet
	if state.HasArtifact("manual-out") {
		t.Error("Manual step artifact should not exist until completed")
	}
}

func TestListWaitingTasks(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	wf := &workflow.Workflow{
		ID: "test",
		Steps: []workflow.Step{
			{Name: "auto1", Handler: "tool", Content: "initial", Output: "out1"},
			{Name: "manual1", Handler: "human", Prompt: "Review data 1", Inputs: []string{"out1"}, Output: "out2"},
			{Name: "manual2", Handler: "human", Prompt: "Review data 2", Inputs: []string{"out1"}, Output: "out3"},
			{Name: "auto2", Handler: "tool", Inputs: []string{"out2", "out3"}, Output: "out4"},
		},
	}

	runName := "test-run"
	CreateRun(wf, runName)

	// Initially no tasks waiting
	tasks, err := ListWaitingTasks(wf, runName)
	if err != nil {
		t.Fatalf("ListWaitingTasks failed: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("Expected 0 waiting tasks initially, got %d", len(tasks))
	}

	// After first tick, auto1 completes
	Tick(wf, runName)

	// After second tick, both manual tasks should be ready
	Tick(wf, runName)

	tasks, err = ListWaitingTasks(wf, runName)
	if err != nil {
		t.Fatalf("ListWaitingTasks failed: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("Expected 2 waiting tasks, got %d", len(tasks))
	}

	// Verify task details
	if tasks[0].Name != "manual1" && tasks[1].Name != "manual1" {
		t.Error("manual1 should be in waiting tasks")
	}
	if tasks[0].Name != "manual2" && tasks[1].Name != "manual2" {
		t.Error("manual2 should be in waiting tasks")
	}

	// Find manual1 task
	var manual1Task WaitingTask
	for _, task := range tasks {
		if task.Name == "manual1" {
			manual1Task = task
			break
		}
	}

	if manual1Task.Prompt != "Review data 1" {
		t.Errorf("Expected prompt 'Review data 1', got '%s'", manual1Task.Prompt)
	}
	if len(manual1Task.Inputs) != 1 || manual1Task.Inputs[0] != "out1" {
		t.Errorf("Expected inputs [out1], got %v", manual1Task.Inputs)
	}
	if manual1Task.Output != "out2" {
		t.Errorf("Expected output 'out2', got '%s'", manual1Task.Output)
	}
}

func TestCompleteTask(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	wf := &workflow.Workflow{
		ID: "test",
		Steps: []workflow.Step{
			{Name: "auto", Handler: "tool", Content: "auto data", Output: "out1"},
			{Name: "manual1", Handler: "human", Prompt: "Task 1", Inputs: []string{"out1"}, Output: "out2"},
			{Name: "manual2", Handler: "human", Prompt: "Task 2", Inputs: []string{"out1"}, Output: "out3"},
			{Name: "final", Handler: "tool", Inputs: []string{"out2", "out3"}, Output: "out4"},
		},
	}

	runName := "test-run"
	CreateRun(wf, runName)

	// Run until both manual tasks are ready
	Tick(wf, runName) // auto completes
	Tick(wf, runName) // manual tasks become ready

	// Complete task at index 0
	err := CompleteTask(wf, runName, 0)
	if err != nil {
		t.Fatalf("CompleteTask failed: %v", err)
	}

	// Verify the task is now succeeded
	state, _ := workflow.LoadState(runName)
	tasks, _ := ListWaitingTasks(wf, runName)

	// One task should still be waiting
	if len(tasks) != 1 {
		t.Errorf("Expected 1 waiting task after completing one, got %d", len(tasks))
	}

	// The completed task should be succeeded
	taskNames := []string{"manual1", "manual2"}
	succeededCount := 0
	for _, name := range taskNames {
		if state.StepStates[name].Status == workflow.StatusSucceeded {
			succeededCount++
			// Should have the artifact
			expectedOutput := "out2"
			if name == "manual2" {
				expectedOutput = "out3"
			}
			if !state.HasArtifact(expectedOutput) {
				t.Errorf("Should have artifact %s after completing %s", expectedOutput, name)
			}
			// Verify artifact content is from auto step
			content, _ := state.ReadArtifact(expectedOutput)
			if content != "auto data" {
				t.Errorf("Artifact %s should contain 'auto data', got '%s'", expectedOutput, content)
			}
		}
	}

	if succeededCount != 1 {
		t.Errorf("Expected 1 succeeded manual task, got %d", succeededCount)
	}
}

func TestCompleteTaskInvalidIndex(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	wf := &workflow.Workflow{
		ID: "test",
		Steps: []workflow.Step{
			{Name: "manual", Handler: "human", Prompt: "Task", Content: "data", Output: "out"},
		},
	}

	runName := "test-run"
	CreateRun(wf, runName)
	Tick(wf, runName) // Make manual task ready

	// Try to complete with invalid index
	err := CompleteTask(wf, runName, 5)
	if err == nil {
		t.Error("Expected error for invalid task index")
	}
}

func TestMixedHandlerWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	wf := &workflow.Workflow{
		ID: "test",
		Steps: []workflow.Step{
			{Name: "fetch", Handler: "tool", Content: "fetched data", Output: "data"},
			{Name: "review", Handler: "human", Prompt: "Review the data", Inputs: []string{"data"}, Output: "reviewed"},
			{Name: "process", Handler: "tool", Inputs: []string{"reviewed"}, Output: "processed"},
		},
	}

	runName := "test-run"
	CreateRun(wf, runName)

	// First tick: fetch runs and completes
	complete, _ := Tick(wf, runName)
	if complete {
		t.Error("Should not be complete after first tick")
	}

	state, _ := workflow.LoadState(runName)
	if state.StepStates["fetch"].Status != workflow.StatusSucceeded {
		t.Error("fetch should be succeeded")
	}

	// Second tick: review becomes ready
	complete, _ = Tick(wf, runName)
	if complete {
		t.Error("Should not be complete with ready task")
	}

	state, _ = workflow.LoadState(runName)
	if state.StepStates["review"].Status != workflow.StatusReady {
		t.Error("review should be ready")
	}
	if state.StepStates["process"].Status != workflow.StatusPending {
		t.Error("process should still be pending")
	}

	// Complete the review task
	CompleteTask(wf, runName, 0)

	// Third tick: process runs
	complete, _ = Tick(wf, runName)
	if !complete {
		t.Error("Should be complete after processing")
	}

	state, _ = workflow.LoadState(runName)
	if state.StepStates["process"].Status != workflow.StatusSucceeded {
		t.Error("process should be succeeded")
	}
}

func TestDefaultHandlerIsTool(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	wf := &workflow.Workflow{
		ID: "test",
		Steps: []workflow.Step{
			{Name: "step", Content: "default content", Output: "out"}, // No handler specified
		},
	}

	runName := "test-run"
	CreateRun(wf, runName)

	// Should auto-execute like a tool handler
	complete, _ := Tick(wf, runName)
	if !complete {
		t.Error("Should be complete - default handler should be tool")
	}

	state, _ := workflow.LoadState(runName)
	if state.StepStates["step"].Status != workflow.StatusSucceeded {
		t.Error("Step with no handler should auto-execute (default to tool)")
	}
}
