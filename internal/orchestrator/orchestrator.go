package orchestrator

import (
	"fmt"
	"sync"

	"composer/internal/workflow"
)

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
	runnableSteps := findRunnableSteps(wf, state)

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
		wg.Add(1)
		go func(s workflow.Step) {
			defer wg.Done()

			// For now, just print that we're running the step
			fmt.Printf("Running step: %s\n", s.Name)
			fmt.Printf("  Description: %s\n", s.Description)
			if len(s.Inputs) > 0 {
				fmt.Printf("  Inputs: %v\n", s.Inputs)
			}
			fmt.Printf("  Output: %s\n", s.Output)
			fmt.Println()

			// Update state with success and add output
			mu.Lock()
			state.StepStates[s.Name] = workflow.StepState{
				Status: workflow.StatusSucceeded,
			}
			state.AddOutput(s.Output)
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
func findRunnableSteps(wf *workflow.Workflow, state *workflow.RunState) []workflow.Step {
	runnable := []workflow.Step{}

	for _, step := range wf.Steps {
		// Skip if step is not pending
		stepState, exists := state.StepStates[step.Name]
		if !exists || stepState.Status != workflow.StatusPending {
			continue
		}

		// Check if all inputs are satisfied
		canRun := true
		for _, input := range step.Inputs {
			if !state.HasOutput(input) {
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
