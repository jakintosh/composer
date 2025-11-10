package column

import (
	"testing"

	"composer/internal/ui/components/ui/button"
	"composer/internal/ui/components/ui/columnheader"
	"gotest.tools/v3/golden"
)

func TestRenderWorkflowColumn(t *testing.T) {
	props := Props{
		Header: columnheader.Props{
			Title: "Workflows",
			Actions: []button.Props{
				{Label: "Add", Class: "button--accent"},
			},
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

	html, err := Render(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "column.golden")
}
