package pages_test

import (
	"testing"

	"composer/internal/ui/pages"
	"composer/internal/ui/views"
	"composer/pkg/ui/components"
	"composer/pkg/ui/testutil"

	"gotest.tools/v3/golden"
)

func TestRenderDashboardPage(t *testing.T) {
	props := pages.DashboardProps{
		Sidebar: components.SidebarProps{
			Title: "Composer",
			Links: []components.SidebarLink{
				{
					Label:  "Dashboard",
					Href:   "/",
					Active: true,
				},
			},
		},
		WorkflowColumn: views.WorkflowColumnProps{
			Title: "Workflows",
		},
		WorkflowModal: views.WorkflowModalProps{
			AddStepButton: components.ButtonProps{
				Label: "Add Step",
			},
		},
		RunColumn: views.RunColumnProps{
			Title: "Runs",
		},
		RunModal: pages.DefaultRunModal(),
		TaskColumn: views.WaitingColumnProps{
			Title: "Tasks",
		},
	}

	html := testutil.Render(t, pages.Dashboard(props))
	golden.Assert(t, html, "dashboard.golden")
}
