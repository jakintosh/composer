package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// LoadWorkflow searches for a workflow file with the given id in the search paths
// and loads it. Returns the workflow and the path where it was found, or an error.
func LoadWorkflow(id string) (*Workflow, string, error) {
	if id == "" {
		return nil, "", fmt.Errorf("workflow id cannot be empty")
	}

	filename := id + ".toml"
	searchPaths := GetWorkflowPaths()

	for _, dir := range searchPaths {
		workflowPath := filepath.Join(dir, filename)

		// Check if file exists
		if _, err := os.Stat(workflowPath); err != nil {
			if os.IsNotExist(err) {
				continue // Try next path
			}
			return nil, "", fmt.Errorf("error checking workflow file %s: %w", workflowPath, err)
		}

		// File exists, try to read and parse it
		data, err := os.ReadFile(workflowPath)
		if err != nil {
			return nil, "", fmt.Errorf("error reading workflow file %s: %w", workflowPath, err)
		}

		var workflow Workflow
		if err := toml.Unmarshal(data, &workflow); err != nil {
			return nil, "", fmt.Errorf("error parsing workflow file %s: %w", workflowPath, err)
		}

		// Set the workflow ID from the filename (without .toml extension)
		workflow.ID = id

		return &workflow, workflowPath, nil
	}

	// Workflow not found in any search path
	return nil, "", fmt.Errorf("workflow '%s' not found in any of the search paths: %v", id, searchPaths)
}

// ListWorkflows returns all workflows found in the search paths
func ListWorkflows() ([]Workflow, error) {
	workflows := []Workflow{}
	seen := make(map[string]bool)

	searchPaths := GetWorkflowPaths()

	for _, dir := range searchPaths {
		// Check if directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		// Read directory entries
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("error reading workflow directory %s: %w", dir, err)
		}

		for _, entry := range entries {
			// Skip directories and non-.toml files
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".toml" {
				continue
			}

			// Extract workflow ID (filename without .toml extension)
			id := entry.Name()[:len(entry.Name())-5]

			// Skip if we've already seen this workflow (earlier paths take precedence)
			if seen[id] {
				continue
			}
			seen[id] = true

			// Load the workflow
			workflow, _, err := LoadWorkflow(id)
			if err != nil {
				return nil, fmt.Errorf("error loading workflow %s: %w", id, err)
			}

			workflows = append(workflows, *workflow)
		}
	}

	return workflows, nil
}

// SaveWorkflow saves a workflow to the local .composer/workflows/ directory
func SaveWorkflow(workflow *Workflow) error {
	if workflow.ID == "" {
		return fmt.Errorf("workflow ID cannot be empty")
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Create .composer/workflows/ directory if it doesn't exist
	workflowDir := filepath.Join(cwd, ".composer", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return fmt.Errorf("failed to create workflow directory: %w", err)
	}

	// Marshal workflow to TOML
	data, err := toml.Marshal(workflow)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow: %w", err)
	}

	// Write to file
	filename := workflow.ID + ".toml"
	workflowPath := filepath.Join(workflowDir, filename)
	if err := os.WriteFile(workflowPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write workflow file: %w", err)
	}

	return nil
}
