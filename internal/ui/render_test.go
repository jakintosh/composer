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
		Sidebar: sidebarViewModel{
			Title: "Composer",
			Links: []sidebarLinkViewModel{
				{Label: "Dashboard", Href: "/", Active: true},
			},
		},
		WorkflowColumn: workflowColumnViewModel{
			Header: columnHeaderViewModel{
				Title: "Workflows",
				Actions: []uiButtonViewModel{
					{
						ID:        "open-workflow-modal",
						Class:     "primary-action",
						Title:     "Create workflow",
						AriaLabel: "Create workflow",
						Type:      "button",
						IconSize:  16,
					},
				},
			},
			Workflows: []workflowViewModel{
				{
					DisplayName: "Example Workflow",
					ID:          "wf-123",
					Description: "Example description",
					Message:     "Example message",
					StepNames:   []string{"Step A", "Step B"},
				},
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
		RunColumn: runColumnViewModel{
			Header: columnHeaderViewModel{Title: "Runs"},
			Runs: []runViewModel{
				{
					DisplayName:  "First Run",
					ID:           "run-1",
					StateLabel:   "ready",
					StateClass:   "state-ready",
					WorkflowName: "Example Workflow",
					Steps: []runStepViewModel{
						{Name: "Step A", Status: "pending", StatusClass: "state-pending"},
						{Name: "Step B", Status: "ready", StatusClass: "state-ready"},
					},
				},
			},
		},
		TaskColumn: waitingTaskColumnViewModel{
			Header: columnHeaderViewModel{Title: "Waiting Tasks"},
			Groups: []waitingTaskGroupViewModel{
				{
					RunID:          "run-1",
					RunDisplayName: "First Run",
					WorkflowName:   "Example Workflow",
					TaskCount:      1,
					Tasks: []waitingTaskViewModel{
						{
							Name:        "Review doc",
							Description: "Look over the generated document",
							Prompt:      "Confirm the document reads well.",
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := renderer.Page(&buf, "pages/dashboard", vm); err != nil {
		t.Fatalf("renderer.Page() error = %v", err)
	}

	if got := buf.String(); !containsAll(got, "Workflow Dashboard", "Example Workflow", "First Run", "run-1", "Waiting Tasks", "Review doc") {
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
