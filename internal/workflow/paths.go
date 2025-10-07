package workflow

import (
	"os"
	"path/filepath"
)

// GetSearchPaths returns the ordered list of directories to search for workflow files.
// The search order is:
// 1. ./.composer/ (current working directory)
// 2. $XDG_DATA_HOME/composer/ (or ~/.local/share/composer/ if XDG_DATA_HOME is not set)
// 3. /etc/composer/ (system-wide)
func GetSearchPaths() []string {
	paths := make([]string, 0, 3)

	// 1. Current directory .composer/
	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(cwd, ".composer"))
	}

	// 2. User-local directory
	var userDataDir string
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		userDataDir = filepath.Join(xdgDataHome, "composer")
	} else if homeDir, err := os.UserHomeDir(); err == nil {
		userDataDir = filepath.Join(homeDir, ".local", "share", "composer")
	}
	if userDataDir != "" {
		paths = append(paths, userDataDir)
	}

	// 3. System-wide directory
	paths = append(paths, "/etc/composer")

	return paths
}

// GetWorkflowPaths returns the ordered list of directories to search for workflow files,
// with the "workflows" subdirectory appended to each search path.
// The search order is:
// 1. ./.composer/workflows/ (current working directory)
// 2. $XDG_DATA_HOME/composer/workflows/ (or ~/.local/share/composer/workflows/)
// 3. /etc/composer/workflows/ (system-wide)
func GetWorkflowPaths() []string {
	basePaths := GetSearchPaths()
	workflowPaths := make([]string, len(basePaths))

	for i, basePath := range basePaths {
		workflowPaths[i] = filepath.Join(basePath, "workflows")
	}

	return workflowPaths
}
