package main

import (
	"fmt"
	"os"

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
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: workflow name is required\n\n")
			printUsage()
			os.Exit(1)
		}
		workflowName := os.Args[2]
		runWorkflow(workflowName)
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
	fmt.Println("  run <workflow-name>    Run a workflow")
}

func runWorkflow(name string) {
	wf, path, err := workflow.LoadWorkflow(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded workflow from: %s\n", path)
	fmt.Printf("Name: %s\n", wf.Name)
	fmt.Printf("Description: %s\n", wf.Description)
	fmt.Printf("Message: %s\n", wf.Message)
}
