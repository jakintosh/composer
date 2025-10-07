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
