package shell

import (
	"html/template"
	"testing"

	"gotest.tools/v3/golden"
)

func TestRenderShell(t *testing.T) {
	props := Props{
		ID:         "example-modal",
		Title:      "Example",
		CloseLabel: "Close example",
		Body:       template.HTML("<p>Body</p>"),
	}

	html, err := Render(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "shell.golden")
}
