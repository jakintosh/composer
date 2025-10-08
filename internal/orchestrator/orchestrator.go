package orchestrator

import (
	"fmt"
	"sync"

	"composer/internal/workflow"
)

// WaitingTask represents a task that is ready for human intervention
type WaitingTask struct {
	Name        string
	Description string
	Prompt      string
	Inputs      []string
	Output      string
}

// CreateRun initializes a new workflow run with the given name
func CreateRun(wf *workflow.Workflow, runName string) error {
	// Create initial state
	state := workflow.NewRunState(wf)

	// Save the initial state
	if err := workflow.SaveState(runName, state); err != nil {
		return fmt.Errorf("failed to save initial state: %w", err)
	}

	return nil
}

// Tick executes one tick of the workflow, running any steps that are ready
func Tick(wf *workflow.Workflow, runName string) (bool, error) {
	// Load current state
	state, err := workflow.LoadState(runName)
	if err != nil {
		return false, fmt.Errorf("failed to load state: %w", err)
	}

	// Check if workflow is already complete
	if state.AllStepsCompleted() {
		return true, nil
	}

	// Find all runnable steps
	runnableSteps := findRunnableSteps(wf, state, runName)

	if len(runnableSteps) == 0 {
		// No steps can run, but workflow isn't complete
		// This could mean we're waiting for something or there's a deadlock
		return false, nil
	}

	// Execute runnable steps in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := []error{}

	for _, step := range runnableSteps {
		// Check if this is a human-handled step
		handler := step.Handler
		if handler == "" {
			handler = "tool" // default to tool
		}

		if handler == "human" {
			// Don't execute, just mark as ready
			mu.Lock()
			state.StepStates[step.Name] = workflow.StepState{
				Status: workflow.StatusReady,
			}
			mu.Unlock()
			fmt.Printf("Step '%s' is ready for human intervention\n", step.Name)
			continue
		}

		// Execute tool steps
		wg.Add(1)
		go func(s workflow.Step) {
			defer wg.Done()

			// Print step execution info
			fmt.Printf("Running step: %s\n", s.Name)
			fmt.Printf("  Description: %s\n", s.Description)
			if len(s.Inputs) > 0 {
				fmt.Printf("  Inputs: %v\n", s.Inputs)
			}
			fmt.Printf("  Output: %s\n", s.Output)
			fmt.Println()

			// Prepare content for the artifact
			var content string
			if len(s.Inputs) > 0 {
				// Load and concatenate input artifacts
				artifacts, err := workflow.ReadArtifacts(runName, s.Inputs)
				if err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf("failed to read input artifacts for %s: %w", s.Name, err))
					mu.Unlock()
					return
				}

				// Concatenate artifacts in order
				for _, inputName := range s.Inputs {
					content += artifacts[inputName]
				}
			} else {
				// Use inline content for steps with no inputs
				content = s.Content
			}

			// Write output artifact
			if err := workflow.WriteArtifact(runName, s.Output, content); err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to write artifact for %s: %w", s.Name, err))
				mu.Unlock()
				return
			}

			// Update state with success
			mu.Lock()
			state.StepStates[s.Name] = workflow.StepState{
				Status: workflow.StatusSucceeded,
			}
			mu.Unlock()
		}(step)
	}

	// Wait for all steps to complete
	wg.Wait()

	// Check for errors
	if len(errors) > 0 {
		// For now, just return the first error
		// In the future, we might want to handle multiple errors differently
		return false, errors[0]
	}

	// Save updated state
	if err := workflow.SaveState(runName, state); err != nil {
		return false, fmt.Errorf("failed to save state: %w", err)
	}

	// Return whether workflow is complete
	return state.AllStepsCompleted(), nil
}

// findRunnableSteps returns all steps that can be run based on current state
func findRunnableSteps(wf *workflow.Workflow, state *workflow.RunState, runName string) []workflow.Step {
	runnable := []workflow.Step{}

	for _, step := range wf.Steps {
		// Only consider pending steps (not ready, succeeded, or failed)
		stepState, exists := state.StepStates[step.Name]
		if !exists || stepState.Status != workflow.StatusPending {
			continue
		}

		// Check if all inputs are satisfied
		canRun := true
		for _, input := range step.Inputs {
			if !workflow.HasArtifact(runName, input) {
				canRun = false
				break
			}
		}

		if canRun {
			runnable = append(runnable, step)
		}
	}

	return runnable
}

// ListWaitingTasks returns all tasks that are ready for human intervention
func ListWaitingTasks(wf *workflow.Workflow, runName string) ([]WaitingTask, error) {
	// Load current state
	state, err := workflow.LoadState(runName)
	if err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	tasks := []WaitingTask{}

	// Find all steps with status "ready"
	for _, step := range wf.Steps {
		stepState, exists := state.StepStates[step.Name]
		if !exists || stepState.Status != workflow.StatusReady {
			continue
		}

		tasks = append(tasks, WaitingTask{
			Name:        step.Name,
			Description: step.Description,
			Prompt:      step.Prompt,
			Inputs:      step.Inputs,
			Output:      step.Output,
		})
	}

	return tasks, nil
}

// CompleteTask marks a ready task as complete and adds its output
func CompleteTask(wf *workflow.Workflow, runName string, taskIndex int) error {
	// Load current state
	state, err := workflow.LoadState(runName)
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// Get list of waiting tasks
	tasks, err := ListWaitingTasks(wf, runName)
	if err != nil {
		return fmt.Errorf("failed to list waiting tasks: %w", err)
	}

	// Validate task index
	if taskIndex < 0 || taskIndex >= len(tasks) {
		return fmt.Errorf("invalid task index: %d (must be between 0 and %d)", taskIndex, len(tasks)-1)
	}

	// Get the task to complete
	task := tasks[taskIndex]

	// Find the step in the workflow to get its inputs and content
	var step *workflow.Step
	for _, s := range wf.Steps {
		if s.Name == task.Name {
			step = &s
			break
		}
	}
	if step == nil {
		return fmt.Errorf("step %s not found in workflow", task.Name)
	}

	// Prepare content for the artifact
	var content string
	if len(step.Inputs) > 0 {
		// Load and concatenate input artifacts
		artifacts, err := workflow.ReadArtifacts(runName, step.Inputs)
		if err != nil {
			return fmt.Errorf("failed to read input artifacts: %w", err)
		}

		// Concatenate artifacts in order
		for _, inputName := range step.Inputs {
			content += artifacts[inputName]
		}
	} else {
		// Use inline content for steps with no inputs
		content = step.Content
	}

	// Write output artifact
	if err := workflow.WriteArtifact(runName, step.Output, content); err != nil {
		return fmt.Errorf("failed to write artifact: %w", err)
	}

	// Mark step as succeeded
	state.StepStates[task.Name] = workflow.StepState{
		Status: workflow.StatusSucceeded,
	}

	// Save state
	if err := workflow.SaveState(runName, state); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	return nil
}
