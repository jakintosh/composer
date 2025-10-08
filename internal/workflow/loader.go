package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// LoadWorkflow searches for a workflow file with the given name in the search paths
// and loads it. Returns the workflow and the path where it was found, or an error.
func LoadWorkflow(name string) (*Workflow, string, error) {
	if name == "" {
		return nil, "", fmt.Errorf("workflow name cannot be empty")
	}

	filename := name + ".toml"
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
		workflow.ID = name

		return &workflow, workflowPath, nil
	}

	// Workflow not found in any search path
	return nil, "", fmt.Errorf("workflow '%s' not found in any of the search paths: %v", name, searchPaths)
}
