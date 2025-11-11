package views

import (
	"testing"

	"composer/pkg/ui/testutil"
	"gotest.tools/v3/golden"
)

func TestRenderWaitingColumn(t *testing.T) {
	props := WaitingColumnProps{
		Title: "Tasks",
		Groups: []WaitingGroup{
			{
				RunID:          "run-a",
				RunDisplayName: "Run A",
				WorkflowName:   "Alpha",
				TaskCount:      1,
				Tasks: []WaitingTask{
					{Name: "Review", Description: "Check"},
				},
			},
		},
	}

	html := testutil.Render(t, WaitingColumn(props))
	golden.Assert(t, html, "waiting_column.golden")
}
