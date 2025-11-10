package column

import (
	"testing"

	"composer/pkg/ui/components/columnheader"
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

	html, err := Render(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "column.golden")
}
