package waiting

import (
	"testing"

	"composer/internal/ui/testutil"
	"gotest.tools/v3/golden"
)

func TestRenderWaitingColumn(t *testing.T) {
	props := ColumnProps{
		Title: "Tasks",
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

	html := testutil.Render(t, Column(props))
	golden.Assert(t, html, "column.golden")
}
