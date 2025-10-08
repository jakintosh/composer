package main

import (
	"fmt"
	"os"

	"composer/internal/orchestrator"
	"composer/internal/workflow"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "Error: both workflow name and run name are required\n\n")
			printUsage()
			os.Exit(1)
		}
		workflowName := os.Args[2]
		runName := os.Args[3]
		runWorkflow(workflowName, runName)
	case "tick":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: run name is required\n\n")
			printUsage()
			os.Exit(1)
		}
		runName := os.Args[2]
		tickWorkflow(runName)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: composer <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  run <workflow-name> <run-name>    Create and start a workflow run")
	fmt.Println("  tick <run-name>                    Execute one tick of a workflow run")
}

func runWorkflow(workflowName, runName string) {
	// Load the workflow
	wf, path, err := workflow.LoadWorkflow(workflowName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded workflow from: %s\n", path)
	fmt.Printf("ID: %s\n", wf.ID)
	if wf.Title != "" {
		fmt.Printf("Title: %s\n", wf.Title)
	}
	fmt.Printf("Description: %s\n", wf.Description)
	if wf.Message != "" {
		fmt.Printf("Message: %s\n", wf.Message)
	}
	fmt.Println()

	// Create the run
	if err := orchestrator.CreateRun(wf, runName); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating run: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created run: %s\n", runName)
	fmt.Println()

	// Execute first tick
	fmt.Println("Executing first tick...")
	fmt.Println()

	complete, err := orchestrator.Tick(wf, runName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing tick: %v\n", err)
		os.Exit(1)
	}

	if complete {
		fmt.Println("Workflow complete!")
	} else {
		fmt.Printf("Tick complete. Run 'composer tick %s' to continue.\n", runName)
	}
}

func tickWorkflow(runName string) {
	// Load the run state to get the workflow name
	state, err := workflow.LoadState(runName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading run state: %v\n", err)
		fmt.Fprintf(os.Stderr, "Make sure the run '%s' exists.\n", runName)
		os.Exit(1)
	}

	// Load the workflow using the name from the state
	wf, _, err := workflow.LoadWorkflow(state.WorkflowName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading workflow '%s': %v\n", state.WorkflowName, err)
		os.Exit(1)
	}

	// Execute tick
	complete, err := orchestrator.Tick(wf, runName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing tick: %v\n", err)
		os.Exit(1)
	}

	if complete {
		fmt.Println("Workflow complete!")
	} else {
		fmt.Printf("Tick complete. Run 'composer tick %s' to continue.\n", runName)
	}
}
