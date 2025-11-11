package views_test

import (
	"testing"

	"composer/internal/ui/views"
	"composer/pkg/ui/testutil"
	"gotest.tools/v3/golden"
)

func TestRenderRunColumn(t *testing.T) {
	props := views.RunColumnProps{
		Title: "Runs",
		Runs: []views.RunView{
			{
				DisplayName:  "Run A",
				ID:           "run-a",
				StateLabel:   "ready",
				StateClass:   "status-badge--ready",
				WorkflowName: "Alpha",
				Steps: []views.RunStep{
					{
						Name:        "first",
						Status:      "pending",
						StatusClass: "status-badge--pending",
					},
				},
			},
		},
	}

	html := testutil.Render(t, views.RunColumn(props))
	golden.Assert(t, html, "run_column.golden")
}

func TestRenderRunModal(t *testing.T) {
	html := testutil.Render(t, views.RunModal(views.RunModalProps{}))
	golden.Assert(t, html, "run_modal.golden")
}
