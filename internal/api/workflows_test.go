package api_test

import (
	"net/http"
	"testing"

	"composer/internal/workflow"
)

// TestGetWorkflows_Empty tests listing workflows when none exist
func TestGetWorkflows_Empty(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	router := setupRouter()

	// get workflows
	var response struct {
		Error *apiError           `json:"error"`
		Data  []workflow.Workflow `json:"data"`
	}
	result := get(router, "/api/workflows", &response)

	// verify result
	err := expectStatus(http.StatusOK, result)
	if err != nil {
		t.Fatalf("%v\n%v", err, response)
	}

	// validate response
	if len(response.Data) != 0 {
		t.Errorf("Expected empty list, got %d workflows", len(response.Data))
	}
}

// TestGetWorkflows_Multiple tests listing multiple workflows
func TestGetWorkflows_Multiple(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// setup
	createWorkflowFixture(t, "workflow1", "First Workflow")
	createWorkflowFixture(t, "workflow2", "Second Workflow")

	router := setupRouter()

	// get workflows
	var response struct {
		Error *apiError           `json:"error"`
		Data  []workflow.Workflow `json:"data"`
	}
	result := get(router, "/api/workflows", &response)

	// verify result
	err := expectStatus(http.StatusOK, result)
	if err != nil {
		t.Fatalf("%v\n%v", err, response)
	}

	// validate response
	if len(response.Data) != 2 {
		t.Fatalf("Expected 2 workflows, got %d", len(response.Data))
	}

	foundFirst := false
	foundSecond := false
	for _, wf := range response.Data {
		if wf.ID == "workflow1" && wf.Title == "First Workflow" {
			foundFirst = true
		}
		if wf.ID == "workflow2" && wf.Title == "Second Workflow" {
			foundSecond = true
		}
	}
	if !foundFirst {
		t.Error("Expected to find workflow1")
	}
	if !foundSecond {
		t.Error("Expected to find workflow2")
	}
}

// TestGetWorkflow_Success tests retrieving a specific workflow
func TestGetWorkflow_Success(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// setup
	createWorkflowFixture(t, "test-workflow", "Test Workflow")

	router := setupRouter()

	// get workflow
	var response struct {
		Error *apiError         `json:"error"`
		Data  workflow.Workflow `json:"data"`
	}
	result := get(router, "/api/workflow/test-workflow", &response)

	// verify result
	err := expectStatus(http.StatusOK, result)
	if err != nil {
		t.Fatalf("%v\n%v", err, response)
	}

	// validate response
	if response.Data.ID != "test-workflow" {
		t.Errorf("Expected ID 'test-workflow', got '%s'", response.Data.ID)
	}
	if response.Data.Title != "Test Workflow" {
		t.Errorf("Expected title 'Test Workflow', got '%s'", response.Data.Title)
	}
}

// TestGetWorkflow_NotFound tests retrieving a non-existent workflow
func TestGetWorkflow_NotFound(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	router := setupRouter()

	// get workflow
	var response struct {
		Error *apiError         `json:"error"`
		Data  workflow.Workflow `json:"data"`
	}
	result := get(router, "/api/workflow/nonexistent", &response)

	// verify result
	err := expectStatus(http.StatusNotFound, result)
	if err != nil {
		t.Fatalf("%v\n%v", err, response)
	}
}

// TestPostWorkflow_Create tests creating a new workflow
func TestPostWorkflow_Create(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	router := setupRouter()

	// setup
	body := `{
		"title": "New Workflow",
		"description": "A new test workflow",
		"message": "Hello",
		"steps": [
			{
				"name": "step1",
				"description": "First step",
				"output": "result1",
				"content": "Step content"
			}
		]
	}`

	// post workflow
	var response struct {
		Error *apiError         `json:"error"`
		Data  workflow.Workflow `json:"data"`
	}
	result := post(router, "/api/workflow/new-workflow", body, &response)

	// verify result
	err := expectStatus(http.StatusOK, result)
	if err != nil {
		t.Fatalf("%v\n%v", err, response)
	}

	// validate response
	if response.Data.ID != "new-workflow" {
		t.Errorf("Expected ID 'new-workflow', got '%s'", response.Data.ID)
	}
	if response.Data.Title != "New Workflow" {
		t.Errorf("Expected title 'New Workflow', got '%s'", response.Data.Title)
	}
	if len(response.Data.Steps) != 1 {
		t.Errorf("Expected 1 step, got %d", len(response.Data.Steps))
	}

	// Verify workflow was actually saved
	wf, _, err2 := workflow.LoadWorkflow("new-workflow")
	if err2 != nil {
		t.Errorf("Failed to load saved workflow: %v", err2)
	}
	if wf.Title != "New Workflow" {
		t.Errorf("Saved workflow has wrong title: %s", wf.Title)
	}
}

// TestPostWorkflow_InvalidJSON tests creating a workflow with malformed JSON
func TestPostWorkflow_InvalidJSON(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	router := setupRouter()

	// setup - malformed JSON
	body := `{
		"title": "Invalid
		missing closing brace
	`

	// post workflow
	var response struct {
		Error *apiError         `json:"error"`
		Data  workflow.Workflow `json:"data"`
	}
	result := post(router, "/api/workflow/invalid", body, &response)

	// verify result
	err := expectStatus(http.StatusBadRequest, result)
	if err != nil {
		t.Fatalf("%v\n%v", err, response)
	}
}

// TestPostWorkflow_Update tests updating an existing workflow
func TestPostWorkflow_Update(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// setup - create initial workflow
	createWorkflowFixture(t, "existing", "Original Title")

	router := setupRouter()

	// Update with new title
	body := `{
		"title": "Updated Title",
		"description": "Updated description",
		"message": "Updated message"
	}`

	// post workflow
	var response struct {
		Error *apiError         `json:"error"`
		Data  workflow.Workflow `json:"data"`
	}
	result := post(router, "/api/workflow/existing", body, &response)

	// verify result
	err := expectStatus(http.StatusOK, result)
	if err != nil {
		t.Fatalf("%v\n%v", err, response)
	}

	// validate response
	if response.Data.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", response.Data.Title)
	}

	// Verify workflow was actually updated
	wf, _, err2 := workflow.LoadWorkflow("existing")
	if err2 != nil {
		t.Errorf("Failed to load updated workflow: %v", err2)
	}
	if wf.Title != "Updated Title" {
		t.Errorf("Updated workflow has wrong title: %s", wf.Title)
	}
}
