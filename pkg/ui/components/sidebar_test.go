package components_test

import (
	"testing"

	"composer/pkg/ui/components"
	"composer/pkg/ui/testutil"

	"gotest.tools/v3/golden"
)

func TestRenderSidebar(t *testing.T) {
	props := components.SidebarProps{
		Title: "Composer",
		Links: []components.SidebarLink{
			{
				Label:  "Dashboard",
				Href:   "/",
				Active: true,
			},
			{
				Label: "Runs",
				Href:  "/runs",
			},
		},
	}

	html := testutil.Render(t, components.Sidebar(props))
	golden.Assert(t, html, "sidebar.golden")
}
