package workflow

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAndReadArtifact(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runName := "test-run"
	artifactName := "test-artifact"
	content := "This is test content\nWith multiple lines"

	// Write artifact
	err := WriteArtifact(runName, artifactName, content)
	if err != nil {
		t.Fatalf("WriteArtifact failed: %v", err)
	}

	// Verify file was created
	artifactPath := filepath.Join(GetArtifactsDir(runName), artifactName)
	if _, err := os.Stat(artifactPath); os.IsNotExist(err) {
		t.Fatalf("Artifact file was not created at %s", artifactPath)
	}

	// Read artifact
	readContent, err := ReadArtifact(runName, artifactName)
	if err != nil {
		t.Fatalf("ReadArtifact failed: %v", err)
	}

	// Verify content matches
	if readContent != content {
		t.Errorf("Content mismatch.\nExpected: %s\nGot: %s", content, readContent)
	}
}

func TestHasArtifact(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runName := "test-run"

	// Should return false for non-existent artifact
	if HasArtifact(runName, "nonexistent") {
		t.Error("HasArtifact should return false for non-existent artifact")
	}

	// Write an artifact
	WriteArtifact(runName, "existing", "content")

	// Should return true for existing artifact
	if !HasArtifact(runName, "existing") {
		t.Error("HasArtifact should return true for existing artifact")
	}

	// Should still return false for different artifact
	if HasArtifact(runName, "other") {
		t.Error("HasArtifact should return false for different artifact")
	}
}

func TestListArtifacts(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runName := "test-run"

	// Should return empty list when no artifacts exist
	artifacts, err := ListArtifacts(runName)
	if err != nil {
		t.Fatalf("ListArtifacts failed: %v", err)
	}
	if len(artifacts) != 0 {
		t.Errorf("Expected 0 artifacts, got %d", len(artifacts))
	}

	// Write some artifacts
	WriteArtifact(runName, "artifact1", "content1")
	WriteArtifact(runName, "artifact2", "content2")
	WriteArtifact(runName, "artifact3", "content3")

	// List should return all artifacts
	artifacts, err = ListArtifacts(runName)
	if err != nil {
		t.Fatalf("ListArtifacts failed: %v", err)
	}

	if len(artifacts) != 3 {
		t.Errorf("Expected 3 artifacts, got %d", len(artifacts))
	}

	// Verify all artifact names are present
	expectedNames := map[string]bool{
		"artifact1": false,
		"artifact2": false,
		"artifact3": false,
	}

	for _, name := range artifacts {
		if _, exists := expectedNames[name]; exists {
			expectedNames[name] = true
		} else {
			t.Errorf("Unexpected artifact name: %s", name)
		}
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("Expected artifact %s not found in list", name)
		}
	}
}

func TestReadArtifacts(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runName := "test-run"

	// Write multiple artifacts
	WriteArtifact(runName, "doc1", "Content of document 1")
	WriteArtifact(runName, "doc2", "Content of document 2")
	WriteArtifact(runName, "doc3", "Content of document 3")

	// Read multiple artifacts
	names := []string{"doc1", "doc3"}
	artifacts, err := ReadArtifacts(runName, names)
	if err != nil {
		t.Fatalf("ReadArtifacts failed: %v", err)
	}

	if len(artifacts) != 2 {
		t.Errorf("Expected 2 artifacts, got %d", len(artifacts))
	}

	// Verify contents
	if artifacts["doc1"] != "Content of document 1" {
		t.Errorf("doc1 content mismatch: got %s", artifacts["doc1"])
	}
	if artifacts["doc3"] != "Content of document 3" {
		t.Errorf("doc3 content mismatch: got %s", artifacts["doc3"])
	}

	// doc2 should not be in the map
	if _, exists := artifacts["doc2"]; exists {
		t.Error("doc2 should not be in the returned artifacts")
	}
}

func TestReadArtifactsNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runName := "test-run"

	// Try to read a non-existent artifact
	names := []string{"nonexistent"}
	_, err := ReadArtifacts(runName, names)
	if err == nil {
		t.Error("ReadArtifacts should return error for non-existent artifact")
	}
}

func TestWriteArtifactCreatesDirectory(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runName := "new-run"

	// Verify artifacts directory doesn't exist yet
	artifactsDir := GetArtifactsDir(runName)
	if _, err := os.Stat(artifactsDir); !os.IsNotExist(err) {
		t.Fatalf("Artifacts directory should not exist yet")
	}

	// Write artifact should create the directory
	err := WriteArtifact(runName, "first", "content")
	if err != nil {
		t.Fatalf("WriteArtifact failed: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(artifactsDir); os.IsNotExist(err) {
		t.Error("WriteArtifact should have created artifacts directory")
	}
}

func TestReadArtifactEmpty(t *testing.T) {
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	runName := "test-run"

	// Write empty artifact
	err := WriteArtifact(runName, "empty", "")
	if err != nil {
		t.Fatalf("WriteArtifact failed: %v", err)
	}

	// Read empty artifact
	content, err := ReadArtifact(runName, "empty")
	if err != nil {
		t.Fatalf("ReadArtifact failed: %v", err)
	}

	if content != "" {
		t.Errorf("Expected empty content, got: %s", content)
	}
}
