package column

import (
	"testing"

	"composer/internal/ui/components/ui/columnheader"
	"gotest.tools/v3/golden"
)

func TestRenderWaitingColumn(t *testing.T) {
	props := Props{
		Header: columnheader.Props{Title: "Tasks"},
		Groups: []Group{
			{
				RunID:          "run-a",
				RunDisplayName: "Run A",
				WorkflowName:   "Alpha",
				TaskCount:      1,
				Tasks: []Task{
					{Name: "Review", Description: "Check"},
				},
			},
		},
	}

	html, err := Render(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "column.golden")
}
