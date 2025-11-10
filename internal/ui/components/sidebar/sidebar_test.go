package sidebar

import (
	"testing"

	"composer/internal/ui/testutil"
	"gotest.tools/v3/golden"
)

func TestRenderSidebar(t *testing.T) {
	props := Props{
		Title: "Composer",
		Links: []Link{
			{Label: "Dashboard", Href: "/", Active: true},
			{Label: "Runs", Href: "/runs"},
		},
	}

	html := testutil.Render(t, Sidebar(props))
	golden.Assert(t, html, "sidebar.golden")
}
