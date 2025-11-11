package components

import (
	"testing"

	"composer/internal/ui/testutil"
	"gotest.tools/v3/golden"
)

func TestRenderSidebar(t *testing.T) {
	props := SidebarProps{
		Title: "Composer",
		Links: []SidebarLink{
			{Label: "Dashboard", Href: "/", Active: true},
			{Label: "Runs", Href: "/runs"},
		},
	}

	html := testutil.Render(t, Sidebar(props))
	golden.Assert(t, html, "sidebar.golden")
}
