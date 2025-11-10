package column

import (
	"testing"

	"composer/internal/ui/components/columnheader"
	"composer/internal/ui/testutil"
	"gotest.tools/v3/golden"
)

func TestRenderRunColumn(t *testing.T) {
	props := Props{
		Header: columnheader.Props{Title: "Runs"},
		Runs: []Run{
			{
				DisplayName:  "Run A",
				ID:           "run-a",
				StateLabel:   "ready",
				StateClass:   "status-badge--ready",
				WorkflowName: "Alpha",
				Steps: []Step{
					{Name: "first", Status: "pending", StatusClass: "status-badge--pending"},
				},
			},
		},
	}

	html := testutil.Render(t, Column(props))
	golden.Assert(t, html, "column.golden")
}
