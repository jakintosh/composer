package workflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetSearchPaths(t *testing.T) {
	// Get search paths
	paths := GetSearchPaths()

	// Should have 3 paths
	if len(paths) != 3 {
		t.Errorf("Expected 3 search paths, got %d", len(paths))
	}

	// First path should be ./.composer in current working directory
	cwd, _ := os.Getwd()
	expectedFirst := filepath.Join(cwd, ".composer")
	if paths[0] != expectedFirst {
		t.Errorf("First path should be %s, got %s", expectedFirst, paths[0])
	}

	// Second path should contain "composer" in user data directory
	if !strings.Contains(paths[1], "composer") {
		t.Errorf("Second path should contain 'composer', got %s", paths[1])
	}

	// Third path should be /etc/composer
	if paths[2] != "/etc/composer" {
		t.Errorf("Third path should be /etc/composer, got %s", paths[2])
	}
}

func TestGetSearchPathsWithXDGDataHome(t *testing.T) {
	// Set XDG_DATA_HOME environment variable
	originalXDG := os.Getenv("XDG_DATA_HOME")
	defer os.Setenv("XDG_DATA_HOME", originalXDG)

	testXDGPath := "/tmp/test-xdg-data"
	os.Setenv("XDG_DATA_HOME", testXDGPath)

	paths := GetSearchPaths()

	// Second path should use XDG_DATA_HOME
	expectedUserPath := filepath.Join(testXDGPath, "composer")
	if paths[1] != expectedUserPath {
		t.Errorf("Expected second path to be %s, got %s", expectedUserPath, paths[1])
	}
}

func TestGetWorkflowPaths(t *testing.T) {
	// Get workflow paths
	paths := GetWorkflowPaths()

	// Should have 3 paths
	if len(paths) != 3 {
		t.Errorf("Expected 3 workflow paths, got %d", len(paths))
	}

	// First path should be ./.composer/workflows in current working directory
	cwd, _ := os.Getwd()
	expectedFirst := filepath.Join(cwd, ".composer", "workflows")
	if paths[0] != expectedFirst {
		t.Errorf("First path should be %s, got %s", expectedFirst, paths[0])
	}

	// Second path should contain "composer/workflows" in user data directory
	if !strings.Contains(paths[1], "composer") || !strings.HasSuffix(paths[1], "workflows") {
		t.Errorf("Second path should contain 'composer' and end with 'workflows', got %s", paths[1])
	}

	// Third path should be /etc/composer/workflows
	if paths[2] != "/etc/composer/workflows" {
		t.Errorf("Third path should be /etc/composer/workflows, got %s", paths[2])
	}
}

func TestGetWorkflowPathsWithXDGDataHome(t *testing.T) {
	// Set XDG_DATA_HOME environment variable
	originalXDG := os.Getenv("XDG_DATA_HOME")
	defer os.Setenv("XDG_DATA_HOME", originalXDG)

	testXDGPath := "/tmp/test-xdg-data"
	os.Setenv("XDG_DATA_HOME", testXDGPath)

	paths := GetWorkflowPaths()

	// Second path should use XDG_DATA_HOME with workflows subdirectory
	expectedUserPath := filepath.Join(testXDGPath, "composer", "workflows")
	if paths[1] != expectedUserPath {
		t.Errorf("Expected second path to be %s, got %s", expectedUserPath, paths[1])
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
