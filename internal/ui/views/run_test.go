package views

import (
	"testing"

	"composer/pkg/ui/testutil"
	"gotest.tools/v3/golden"
)

func TestRenderRunColumn(t *testing.T) {
	props := RunColumnProps{
		Title: "Runs",
		Runs: []RunView{
			{
				DisplayName:  "Run A",
				ID:           "run-a",
				StateLabel:   "ready",
				StateClass:   "status-badge--ready",
				WorkflowName: "Alpha",
				Steps: []RunStep{
					{Name: "first", Status: "pending", StatusClass: "status-badge--pending"},
				},
			},
		},
	}

	html := testutil.Render(t, RunColumn(props))
	golden.Assert(t, html, "run_column.golden")
}

func TestRenderRunModal(t *testing.T) {
	html := testutil.Render(t, RunModal(RunModalProps{}))
	golden.Assert(t, html, "run_modal.golden")
}
