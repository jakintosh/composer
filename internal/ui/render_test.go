package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestRendererLoadsDashboardTemplate(t *testing.T) {
	renderer, err := NewRenderer(nil)
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}

	if renderer.templates.Lookup("layouts/base") == nil {
		t.Fatalf("expected layouts/base template to be parsed")
	}

	vm := dashboardViewModel{
		WorkflowColumn: workflowColumnViewModel{
			Workflows: []workflowViewModel{
				{
					DisplayName: "Example Workflow",
					ID:          "wf-123",
					Title:       "Example Workflow",
					Description: "Example description",
					Message:     "Example message",
					StepNames:   []string{"Step A", "Step B"},
				},
			},
			CreateButton: uiButtonViewModel{
				ID:        "open-workflow-modal",
				Class:     "primary-action",
				Title:     "Create workflow",
				AriaLabel: "Create workflow",
				Type:      "button",
				IconSize:  16,
			},
		},
		WorkflowModal: workflowModalViewModel{
			AddStepButton: uiButtonViewModel{
				ID:       "add-workflow-step",
				Class:    "add-step-button",
				Label:    "Add Step",
				Type:     "button",
				IconSize: 16,
			},
		},
		Runs: []runViewModel{
			{
				Name:         "run-1",
				StateLabel:   "ready",
				StateClass:   "state-ready",
				WorkflowName: "Example Workflow",
				Steps: []runStepViewModel{
					{Name: "Step A", Status: "pending", StatusClass: "state-pending"},
					{Name: "Step B", Status: "ready", StatusClass: "state-ready"},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := renderer.Page(&buf, "pages/dashboard", vm); err != nil {
		t.Fatalf("renderer.Page() error = %v", err)
	}

	if got := buf.String(); !containsAll(got, "Workflow Dashboard", "Example Workflow", "run-1") {
		t.Fatalf("rendered output missing expected content: %q", got)
	}
}

func containsAll(haystack string, needles ...string) bool {
	for _, needle := range needles {
		if !strings.Contains(haystack, needle) {
			return false
		}
	}
	return true
}
