package workflow

import (
	"fmt"
	"os"
	"path/filepath"
)

// HasArtifact checks if an artifact with the given name exists for a run
func HasArtifact(runName, name string) bool {
	artifactPath := filepath.Join(GetArtifactsDir(runName), name)
	_, err := os.Stat(artifactPath)
	return err == nil
}

// ListArtifacts returns a list of all artifact names in a run's artifacts directory
func ListArtifacts(runName string) ([]string, error) {
	artifactsDir := GetArtifactsDir(runName)

	// Check if artifacts directory exists
	if _, err := os.Stat(artifactsDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(artifactsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read artifacts directory: %w", err)
	}

	artifacts := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			artifacts = append(artifacts, entry.Name())
		}
	}

	return artifacts, nil
}

// WriteArtifact writes content to an artifact file
func WriteArtifact(runName, name, content string) error {
	artifactsDir := GetArtifactsDir(runName)

	// Create artifacts directory if it doesn't exist
	if err := os.MkdirAll(artifactsDir, 0755); err != nil {
		return fmt.Errorf("failed to create artifacts directory: %w", err)
	}

	artifactPath := filepath.Join(artifactsDir, name)
	if err := os.WriteFile(artifactPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write artifact %s: %w", name, err)
	}

	return nil
}

// ReadArtifact reads the content of a single artifact
func ReadArtifact(runName, name string) (string, error) {
	artifactPath := filepath.Join(GetArtifactsDir(runName), name)

	data, err := os.ReadFile(artifactPath)
	if err != nil {
		return "", fmt.Errorf("failed to read artifact %s: %w", name, err)
	}

	return string(data), nil
}

// ReadArtifacts reads multiple artifacts and returns a map of name to content
func ReadArtifacts(runName string, names []string) (map[string]string, error) {
	artifacts := make(map[string]string)

	for _, name := range names {
		content, err := ReadArtifact(runName, name)
		if err != nil {
			return nil, err
		}
		artifacts[name] = content
	}

	return artifacts, nil
}
