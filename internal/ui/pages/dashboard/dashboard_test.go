package dashboard

import (
	"path/filepath"
	"testing"

	"composer/internal/ui/components/navigation/sidebar"
	"composer/internal/ui/components/run/column"
	"composer/internal/ui/components/ui/button"
	"composer/internal/ui/components/ui/columnheader"
	waitingcolumn "composer/internal/ui/components/waiting/column"
	workflowcolumn "composer/internal/ui/components/workflow/column"
	workflowmodal "composer/internal/ui/components/workflow/modal"
	"gotest.tools/v3/golden"
)

func TestRenderDashboardPage(t *testing.T) {
	props := Props{
		Sidebar: sidebar.Props{
			Title: "Composer",
			Links: []sidebar.Link{{Label: "Dashboard", Href: "/", Active: true}},
		},
		WorkflowColumn: workflowcolumn.Props{
			Header: columnheader.Props{Title: "Workflows"},
		},
		WorkflowModal: workflowmodal.Props{
			AddStepButton: button.Props{Label: "Add Step"},
		},
		RunColumn:  column.Props{Header: columnheader.Props{Title: "Runs"}},
		RunModal:   DefaultRunModal(),
		TaskColumn: waitingcolumn.Props{Header: columnheader.Props{Title: "Tasks"}},
	}

	html, err := HTML(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), filepath.Join("fixtures", "dashboard.golden"))
}
