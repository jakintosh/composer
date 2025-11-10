package dashboard

import (
	"path/filepath"
	"testing"

	"composer/internal/ui/components/button"
	"composer/internal/ui/components/columnheader"
	column "composer/internal/ui/components/run_column"
	sidebar "composer/internal/ui/components/sidebar"
	waitingcolumn "composer/internal/ui/components/waiting_column"
	workflowcolumn "composer/internal/ui/components/workflow_column"
	workflowmodal "composer/internal/ui/components/workflow_modal"
	"composer/internal/ui/testutil"
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

	html := testutil.Render(t, Page(props))
	golden.Assert(t, html, filepath.Join("fixtures", "dashboard.golden"))
}
