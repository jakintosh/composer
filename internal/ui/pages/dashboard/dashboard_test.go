package dashboard

import (
	"path/filepath"
	"testing"

	"composer/internal/ui/components/button"
	runcomponent "composer/internal/ui/components/run"
	sidebar "composer/internal/ui/components/sidebar"
	waitingcolumn "composer/internal/ui/components/waiting_column"
	workflowcomponent "composer/internal/ui/components/workflow"
	"composer/internal/ui/testutil"
	"gotest.tools/v3/golden"
)

func TestRenderDashboardPage(t *testing.T) {
	props := Props{
		Sidebar: sidebar.Props{
			Title: "Composer",
			Links: []sidebar.Link{{Label: "Dashboard", Href: "/", Active: true}},
		},
		WorkflowColumn: workflowcomponent.ColumnProps{
			Title: "Workflows",
		},
		WorkflowModal: workflowcomponent.ModalProps{
			AddStepButton: button.Props{Label: "Add Step"},
		},
		RunColumn:  runcomponent.ColumnProps{Title: "Runs"},
		RunModal:   DefaultRunModal(),
		TaskColumn: waitingcolumn.Props{Title: "Tasks"},
	}

	html := testutil.Render(t, Page(props))
	golden.Assert(t, html, filepath.Join("fixtures", "dashboard.golden"))
}
