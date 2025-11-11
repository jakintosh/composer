package views_test

import (
	"testing"

	"composer/internal/ui/views"
	"composer/pkg/ui/components"
	"composer/pkg/ui/testutil"

	"gotest.tools/v3/golden"
)

func TestRenderWorkflowColumn(t *testing.T) {
	props := views.WorkflowColumnProps{
		Title: "Workflows",
		Actions: []components.ButtonProps{
			{
				Label: "Add",
				Class: "button--accent",
			},
		},
		Workflows: []views.WorkflowView{
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

	html := testutil.Render(t, views.WorkflowColumn(props))
	golden.Assert(t, html, "workflow_column.golden")
}

func TestRenderWorkflowModal(t *testing.T) {
	props := views.WorkflowModalProps{
		AddStepButton: components.ButtonProps{
			Label: "Add Step",
		},
	}

	html := testutil.Render(t, views.WorkflowModal(props))
	golden.Assert(t, html, "workflow_modal.golden")
}
