package workflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadWorkflow_ValidWorkflow(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create a test workflow file
	workflowContent := `
display_name = "Test Workflow"
description = "A test workflow"
message = "Hello, World!"
`
	workflowPath := filepath.Join(tmpDir, "test.toml")
	if err := os.WriteFile(workflowPath, []byte(workflowContent), 0644); err != nil {
		t.Fatalf("Failed to create test workflow file: %v", err)
	}

	// Change to temp directory so ./.composer is our temp dir
	originalCwd, _ := os.Getwd()
	defer os.Chdir(originalCwd)
	os.Chdir(tmpDir)

	// Create .composer/workflows directory and copy workflow there
	composerDir := filepath.Join(tmpDir, ".composer", "workflows")
	os.MkdirAll(composerDir, 0755)
	composerWorkflowPath := filepath.Join(composerDir, "test.toml")
	os.WriteFile(composerWorkflowPath, []byte(workflowContent), 0644)

	// Load the workflow
	workflow, foundPath, err := LoadWorkflow("test")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if workflow == nil {
		t.Fatal("Expected workflow, got nil")
	}

	if workflow.ID != "test" {
		t.Errorf("Expected ID 'test', got '%s'", workflow.ID)
	}

	if workflow.DisplayName != "Test Workflow" {
		t.Errorf("Expected display name 'Test Workflow', got '%s'", workflow.DisplayName)
	}

	if workflow.Description != "A test workflow" {
		t.Errorf("Expected description 'A test workflow', got '%s'", workflow.Description)
	}

	if workflow.Message != "Hello, World!" {
		t.Errorf("Expected message 'Hello, World!', got '%s'", workflow.Message)
	}

	if !strings.HasSuffix(foundPath, "test.toml") {
		t.Errorf("Expected path to end with 'test.toml', got '%s'", foundPath)
	}
}

func TestLoadWorkflow_NotFound(t *testing.T) {
	// Change to a temp directory where no workflows exist
	tmpDir := t.TempDir()
	originalCwd, _ := os.Getwd()
	defer os.Chdir(originalCwd)
	os.Chdir(tmpDir)

	// Try to load a non-existent workflow
	workflow, _, err := LoadWorkflow("nonexistent")
	if err == nil {
		t.Fatal("Expected error for non-existent workflow, got nil")
	}

	if workflow != nil {
		t.Errorf("Expected nil workflow, got %v", workflow)
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

func TestLoadWorkflow_InvalidTOML(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create an invalid TOML file
	invalidContent := `
title = "test
this is not valid TOML
`
	composerDir := filepath.Join(tmpDir, ".composer", "workflows")
	os.MkdirAll(composerDir, 0755)
	workflowPath := filepath.Join(composerDir, "invalid.toml")
	if err := os.WriteFile(workflowPath, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("Failed to create test workflow file: %v", err)
	}

	// Change to temp directory
	originalCwd, _ := os.Getwd()
	defer os.Chdir(originalCwd)
	os.Chdir(tmpDir)

	// Try to load the invalid workflow
	workflow, _, err := LoadWorkflow("invalid")
	if err == nil {
		t.Fatal("Expected error for invalid TOML, got nil")
	}

	if workflow != nil {
		t.Errorf("Expected nil workflow, got %v", workflow)
	}

	if !strings.Contains(err.Error(), "error parsing") {
		t.Errorf("Expected 'error parsing' in error message, got: %v", err)
	}
}

func TestLoadWorkflow_EmptyID(t *testing.T) {
	workflow, _, err := LoadWorkflow("")
	if err == nil {
		t.Fatal("Expected error for empty workflow id, got nil")
	}

	if workflow != nil {
		t.Errorf("Expected nil workflow, got %v", workflow)
	}

	if !strings.Contains(err.Error(), "cannot be empty") {
		t.Errorf("Expected 'cannot be empty' error, got: %v", err)
	}
}
