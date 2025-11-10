package columnheader

import (
	"testing"

	"composer/internal/ui/components/ui/button"
	"gotest.tools/v3/golden"
)

func TestRenderColumnHeader(t *testing.T) {
	props := Props{
		Title: "Workflows",
		Actions: []button.Props{
			{
				ID:        "create",
				Class:     "button--accent",
				Title:     "Create workflow",
				AriaLabel: "Create workflow",
				Label:     "Create",
				Type:      "button",
			},
		},
	}

	html, err := Render(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "header.golden")
}
