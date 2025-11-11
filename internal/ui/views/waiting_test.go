package views_test

import (
	"testing"

	"composer/internal/ui/views"
	"composer/pkg/ui/testutil"

	"gotest.tools/v3/golden"
)

func TestRenderWaitingColumn(t *testing.T) {
	props := views.WaitingColumnProps{
		Title: "Tasks",
		Groups: []views.WaitingGroup{
			{
				RunID:          "run-a",
				RunDisplayName: "Run A",
				WorkflowName:   "Alpha",
				TaskCount:      1,
				Tasks: []views.WaitingTask{
					{
						Name:        "Review",
						Description: "Check",
					},
				},
			},
		},
	}

	html := testutil.Render(t, views.WaitingColumn(props))
	golden.Assert(t, html, "waiting_column.golden")
}
