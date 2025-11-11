package dashboard

import (
	"path/filepath"
	"testing"

	"composer/pkg/ui/components"
	"composer/internal/ui/testutil"
	"composer/internal/ui/views"
	"gotest.tools/v3/golden"
)

func TestRenderDashboardPage(t *testing.T) {
	props := Props{
		Sidebar: components.SidebarProps{
			Title: "Composer",
			Links: []components.SidebarLink{{Label: "Dashboard", Href: "/", Active: true}},
		},
		WorkflowColumn: views.WorkflowColumnProps{
			Title: "Workflows",
		},
		WorkflowModal: views.WorkflowModalProps{
			AddStepButton: components.ButtonProps{Label: "Add Step"},
		},
		RunColumn:  views.RunColumnProps{Title: "Runs"},
		RunModal:   DefaultRunModal(),
		TaskColumn: views.WaitingColumnProps{Title: "Tasks"},
	}

	html := testutil.Render(t, Page(props))
	golden.Assert(t, html, filepath.Join("fixtures", "dashboard.golden"))
}
