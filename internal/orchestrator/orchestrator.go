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

// CreateRun initializes a new workflow run with the given id and display name
func CreateRun(wf *workflow.Workflow, runID string, displayName string) error {
	// Create initial state
	state := workflow.NewRunState(wf, runID, displayName)

	// Save the initial state
	if err := state.Save(); err != nil {
		return fmt.Errorf("failed to save initial state: %w", err)
	}

	return nil
}

// Tick executes one tick of the workflow, running any steps that are ready
func Tick(wf *workflow.Workflow, runID string) (bool, error) {
	// Load current state
	state, err := workflow.LoadState(runID)
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
				artifacts, err := state.ReadArtifacts(s.Inputs)
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
			if err := state.WriteArtifact(s.Output, content); err != nil {
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
	if err := state.Save(); err != nil {
		return false, fmt.Errorf("failed to save state: %w", err)
	}

	// Return whether workflow is complete
	return state.AllStepsCompleted(), nil
}

// findRunnableSteps returns all steps that can be run based on current state
func findRunnableSteps(wf *workflow.Workflow, state *workflow.RunState) []workflow.Step {
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
			if !state.HasArtifact(input) {
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
func ListWaitingTasks(wf *workflow.Workflow, runID string) ([]WaitingTask, error) {
	// Load current state
	state, err := workflow.LoadState(runID)
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

// ListWaitingTasksByRun returns waiting tasks grouped by run name.
func ListWaitingTasksByRun(runs []workflow.RunState) (map[string][]WaitingTask, error) {
	tasksByRun := make(map[string][]WaitingTask, len(runs))
	workflowCache := make(map[string]*workflow.Workflow)

	for _, run := range runs {
		if run.ID == "" || run.WorkflowName == "" {
			continue
		}

		wf, ok := workflowCache[run.WorkflowName]
		if !ok {
			var err error
			wf, _, err = workflow.LoadWorkflow(run.WorkflowName)
			if err != nil {
				return nil, fmt.Errorf("load workflow '%s' for run '%s': %w", run.WorkflowName, run.ID, err)
			}
			workflowCache[run.WorkflowName] = wf
		}

		tasks, err := ListWaitingTasks(wf, run.ID)
		if err != nil {
			return nil, fmt.Errorf("list waiting tasks for run '%s': %w", run.ID, err)
		}
		tasksByRun[run.ID] = tasks
	}

	return tasksByRun, nil
}

// CompleteTask marks a ready task as complete and adds its output
func CompleteTask(wf *workflow.Workflow, runID string, taskIndex int) error {
	// Load current state
	state, err := workflow.LoadState(runID)
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// Get list of waiting tasks
	tasks, err := ListWaitingTasks(wf, runID)
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
		artifacts, err := state.ReadArtifacts(step.Inputs)
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
	if err := state.WriteArtifact(step.Output, content); err != nil {
		return fmt.Errorf("failed to write artifact: %w", err)
	}

	// Mark step as succeeded
	state.StepStates[task.Name] = workflow.StepState{
		Status: workflow.StatusSucceeded,
	}

	// Save state
	if err := state.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	return nil
}
