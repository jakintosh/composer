package workflow

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestGetSearchPaths(t *testing.T) {
	// Get search paths
	paths := GetSearchPaths()

	// Should have at least 1 path (system path is always added)
	if len(paths) < 1 {
		t.Fatalf("Expected at least 1 search path, got %d", len(paths))
	}

	// Last path should always be /etc/composer
	if paths[len(paths)-1] != "/etc/composer" {
		t.Errorf("Last path should be /etc/composer, got %s", paths[len(paths)-1])
	}

	// If we can get current working directory, first path should be ./.composer
	if cwd, err := os.Getwd(); err == nil {
		expectedFirst := filepath.Join(cwd, ".composer")
		if paths[0] != expectedFirst {
			t.Errorf("First path should be %s, got %s", expectedFirst, paths[0])
		}
	}

	// Should contain a path with "composer" (user data directory)
	hasComposerPath := false
	for _, path := range paths {
		if strings.Contains(path, "composer") {
			hasComposerPath = true
			break
		}
	}
	if !hasComposerPath {
		t.Error("Expected at least one path to contain 'composer'")
	}
}

func TestGetSearchPathsWithXDGDataHome(t *testing.T) {
	// Set XDG_DATA_HOME environment variable
	originalXDG := os.Getenv("XDG_DATA_HOME")
	defer os.Setenv("XDG_DATA_HOME", originalXDG)

	testXDGPath := "/tmp/test-xdg-data"
	os.Setenv("XDG_DATA_HOME", testXDGPath)

	paths := GetSearchPaths()

	// Should contain a path using XDG_DATA_HOME
	expectedUserPath := filepath.Join(testXDGPath, "composer")
	foundXDGPath := slices.Contains(paths, expectedUserPath)
	if !foundXDGPath {
		t.Errorf("Expected to find XDG path %s in paths %v", expectedUserPath, paths)
	}
}

func TestGetWorkflowPaths(t *testing.T) {
	// Get workflow paths
	paths := GetWorkflowPaths()

	// Should have at least 1 path (system path is always added)
	if len(paths) < 1 {
		t.Fatalf("Expected at least 1 workflow path, got %d", len(paths))
	}

	// Last path should always be /etc/composer/workflows
	if paths[len(paths)-1] != "/etc/composer/workflows" {
		t.Errorf("Last path should be /etc/composer/workflows, got %s", paths[len(paths)-1])
	}

	// If we can get current working directory, first path should be ./.composer/workflows
	if cwd, err := os.Getwd(); err == nil {
		expectedFirst := filepath.Join(cwd, ".composer", "workflows")
		if paths[0] != expectedFirst {
			t.Errorf("First path should be %s, got %s", expectedFirst, paths[0])
		}
	}

	// Should contain a path with "composer" and ending with "workflows"
	hasComposerWorkflowsPath := false
	for _, path := range paths {
		if strings.Contains(path, "composer") && strings.HasSuffix(path, "workflows") {
			hasComposerWorkflowsPath = true
			break
		}
	}
	if !hasComposerWorkflowsPath {
		t.Error("Expected at least one path to contain 'composer' and end with 'workflows'")
	}
}

func TestGetWorkflowPathsWithXDGDataHome(t *testing.T) {
	// Set XDG_DATA_HOME environment variable
	originalXDG := os.Getenv("XDG_DATA_HOME")
	defer os.Setenv("XDG_DATA_HOME", originalXDG)

	testXDGPath := "/tmp/test-xdg-data"
	os.Setenv("XDG_DATA_HOME", testXDGPath)

	paths := GetWorkflowPaths()

	// Should contain a path using XDG_DATA_HOME with workflows subdirectory
	expectedUserPath := filepath.Join(testXDGPath, "composer", "workflows")
	foundXDGPath := slices.Contains(paths, expectedUserPath)
	if !foundXDGPath {
		t.Errorf("Expected to find XDG path %s in paths %v", expectedUserPath, paths)
	}
}

func TestGetRunsDir(t *testing.T) {
	runsDir := GetRunsDir()

	// Should be ./.composer/runs in current working directory
	cwd, _ := os.Getwd()
	expected := filepath.Join(cwd, ".composer", "runs")
	if runsDir != expected {
		t.Errorf("Expected runs dir to be %s, got %s", expected, runsDir)
	}
}

func TestGetRunDir(t *testing.T) {
	runName := "test-run"
	runDir := GetRunDir(runName)

	// Should be ./.composer/runs/test-run in current working directory
	cwd, _ := os.Getwd()
	expected := filepath.Join(cwd, ".composer", "runs", runName)
	if runDir != expected {
		t.Errorf("Expected run dir to be %s, got %s", expected, runDir)
	}
}

func TestGetRunDirWithDifferentNames(t *testing.T) {
	tests := []struct {
		name    string
		runName string
	}{
		{"simple name", "my-run"},
		{"with numbers", "run-123"},
		{"with underscores", "my_test_run"},
		{"with dots", "run.v1.0"},
	}

	cwd, _ := os.Getwd()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runDir := GetRunDir(tt.runName)
			expected := filepath.Join(cwd, ".composer", "runs", tt.runName)
			if runDir != expected {
				t.Errorf("Expected run dir to be %s, got %s", expected, runDir)
			}
		})
	}
}
