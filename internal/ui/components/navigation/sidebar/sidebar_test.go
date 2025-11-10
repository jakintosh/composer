package sidebar

import (
	"testing"

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

	html, err := Render(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "sidebar.golden")
}
