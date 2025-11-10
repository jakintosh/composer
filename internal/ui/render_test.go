package ui

import (
	"bytes"
	"strings"
	"testing"

	"composer/internal/ui/view"
)

func TestRendererLoadsDashboardTemplate(t *testing.T) {
	server, err := Init(ModeProduction)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	renderer := server.Renderer()

	if renderer.templates.Lookup("layouts/base") == nil {
		t.Fatalf("expected layouts/base template to be parsed")
	}

	vm := view.Dashboard{
		Sidebar: view.Sidebar{
			Title: "Composer",
			Links: []view.SidebarLink{
				{Label: "Dashboard", Href: "/", Active: true},
			},
		},
		WorkflowColumn: view.WorkflowColumn{
			Header: view.ColumnHeader{
				Title: "Workflows",
				Actions: []view.Button{
					{
						ID:        "open-workflow-modal",
						Class:     "button--accent button--icon",
						Title:     "Create workflow",
						AriaLabel: "Create workflow",
						Type:      "button",
						IconSize:  16,
					},
				},
			},
			Workflows: []view.Workflow{
				{
					DisplayName: "Example Workflow",
					ID:          "wf-123",
					Description: "Example description",
					Message:     "Example message",
					StepNames:   []string{"Step A", "Step B"},
				},
			},
		},
		WorkflowModal: view.WorkflowModal{
			AddStepButton: view.Button{
				ID:       "add-workflow-step",
				Class:    "button--outline button--sm",
				Label:    "Add Step",
				Type:     "button",
				IconSize: 16,
			},
		},
		RunColumn: view.RunColumn{
			Header: view.ColumnHeader{Title: "Runs"},
			Runs: []view.Run{
				{
					DisplayName:  "First Run",
					ID:           "run-1",
					StateLabel:   "ready",
					StateClass:   "status-badge--ready",
					WorkflowName: "Example Workflow",
					Steps: []view.RunStep{
						{Name: "Step A", Status: "pending", StatusClass: "status-badge--pending"},
						{Name: "Step B", Status: "ready", StatusClass: "status-badge--ready"},
					},
				},
			},
		},
		TaskColumn: view.WaitingColumn{
			Header: view.ColumnHeader{Title: "Waiting Tasks"},
			Groups: []view.WaitingGroup{
				{
					RunID:          "run-1",
					RunDisplayName: "First Run",
					WorkflowName:   "Example Workflow",
					TaskCount:      1,
					Tasks: []view.WaitingTask{
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
	if err := renderer.Page(&buf, view.DashboardTemplate, vm); err != nil {
		t.Fatalf("renderer.Page() error = %v", err)
	}

	if got := buf.String(); !containsAll(got, "Workflow Dashboard", "Example Workflow", "First Run", "run-1", "Waiting Tasks", "Review doc", `data-run-id="run-1"`) {
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

func TestInitDevelopmentModeCreatesReloadingRenderer(t *testing.T) {
	server, err := Init(ModeDevelopment)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	renderer := server.Renderer()
	if !renderer.dev {
		t.Fatalf("expected renderer.dev to be true in development mode")
	}
}
