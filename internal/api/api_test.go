package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"composer/internal/api"
	"composer/internal/workflow"
)

// httpResult captures the status code and any decoding error from an HTTP request
type httpResult struct {
	Code  int
	Error error
}

// apiResponse matches the standardized APIResponse wrapper
type apiResponse struct {
	Error *apiError `json:"error"`
	Data  any       `json:"data"`
}

// apiError matches the APIError structure
type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// setupRouter creates a fresh router with all API routes registered
func setupRouter() *http.ServeMux {
	return api.BuildRouter()
}

// setupTestEnv creates a temporary directory and changes to it for testing.
// It returns a cleanup function that should be deferred.
func setupTestEnv(t *testing.T) func() {
	t.Helper()

	// Create temporary directory
	tmpDir := t.TempDir()

	// Save current directory
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create .composer directories
	if err := os.MkdirAll(filepath.Join(tmpDir, ".composer", "workflows"), 0755); err != nil {
		t.Fatalf("Failed to create workflows directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, ".composer", "runs"), 0755); err != nil {
		t.Fatalf("Failed to create runs directory: %v", err)
	}

	// Return cleanup function
	return func() {
		os.Chdir(originalCwd)
	}
}

// expectStatus checks if the HTTP status code matches the expected value
func expectStatus(
	code int,
	result httpResult,
) error {
	if result.Code == code {
		return nil
	}
	return fmt.Errorf("expected status %d, got %d: %v", code, result.Code, result.Error)
}

// get performs a GET request and decodes the response into the provided destination
func get(
	router *http.ServeMux,
	url string,
	response any,
) httpResult {
	req := httptest.NewRequest("GET", url, nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	// Decode response if destination provided
	if response != nil && res.Body.Len() > 0 {
		if err := json.Unmarshal(res.Body.Bytes(), response); err != nil {
			return httpResult{
				Code:  res.Code,
				Error: fmt.Errorf("failed to decode JSON: %v\n%s", err, res.Body.String()),
			}
		}
	}

	return httpResult{Code: res.Code, Error: nil}
}

// post performs a POST request with the given body and decodes the response
func post(
	router *http.ServeMux,
	url string,
	body string,
	response any,
) httpResult {
	req := httptest.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	// Decode response if destination provided
	if response != nil && res.Body.Len() > 0 {
		if err := json.Unmarshal(res.Body.Bytes(), response); err != nil {
			return httpResult{
				Code:  res.Code,
				Error: fmt.Errorf("failed to decode JSON: %v\n%s", err, res.Body.String()),
			}
		}
	}

	return httpResult{Code: res.Code, Error: nil}
}

// createWorkflowFixture creates a test workflow file in .composer/workflows/
func createWorkflowFixture(t *testing.T, id, displayName string) {
	t.Helper()

	content := fmt.Sprintf(`display_name = "%s"
description = "A test workflow"
message = "Test message"

[[steps]]
name = "step1"
description = "First step"
output = "result1"
content = "Step 1 content"
`, displayName)

	workflowDir := filepath.Join(".composer", "workflows")
	workflowPath := filepath.Join(workflowDir, id+".toml")

	if err := os.WriteFile(workflowPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create workflow fixture: %v", err)
	}
}

// createRunFixture creates a test run state in .composer/runs/
func createRunFixture(t *testing.T, runID, workflowID string) {
	t.Helper()

	// Load the workflow to create proper state
	wf, _, err := workflow.LoadWorkflow(workflowID)
	if err != nil {
		t.Fatalf("Failed to load workflow for run fixture: %v", err)
	}

	// Create run state
	state := workflow.NewRunState(wf, runID, runID)
	if err := state.Save(); err != nil {
		t.Fatalf("Failed to save run fixture: %v", err)
	}
}
