package views

import (
	"testing"

	"composer/pkg/ui/components"
	"composer/pkg/ui/testutil"
	"gotest.tools/v3/golden"
)

func TestRenderWorkflowColumn(t *testing.T) {
	props := WorkflowColumnProps{
		Title: "Workflows",
		Actions: []components.ButtonProps{
			{Label: "Add", Class: "button--accent"},
		},
		Workflows: []WorkflowView{
			{
				DisplayName: "Alpha",
				ID:          "alpha",
				Description: "First",
				Message:     "Hello",
				StepNames:   []string{"first", "second"},
			},
			{
				DisplayName: "Beta",
			},
		},
	}

	html := testutil.Render(t, WorkflowColumn(props))
	golden.Assert(t, html, "workflow_column.golden")
}

func TestRenderWorkflowModal(t *testing.T) {
	props := WorkflowModalProps{AddStepButton: components.ButtonProps{Label: "Add Step"}}

	html := testutil.Render(t, WorkflowModal(props))
	golden.Assert(t, html, "workflow_modal.golden")
}
