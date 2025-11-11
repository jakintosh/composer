package workflow

import (
	"testing"

	"composer/internal/ui/components/button"
	"composer/internal/ui/testutil"
	"gotest.tools/v3/golden"
)

func TestRenderWorkflowColumn(t *testing.T) {
	props := ColumnProps{
		Title: "Workflows",
		Actions: []button.Props{
			{Label: "Add", Class: "button--accent"},
		},
		Workflows: []Workflow{
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

	html := testutil.Render(t, Column(props))
	golden.Assert(t, html, "column.golden")
}
