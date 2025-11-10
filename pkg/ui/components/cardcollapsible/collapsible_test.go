package collapsible

import (
	"html/template"
	"testing"

	"gotest.tools/v3/golden"
)

func TestRenderCollapsibleCard(t *testing.T) {
	props := Props{
		Title: "Example",
		SummaryItems: []template.HTML{
			template.HTML("<span>Extra</span>"),
		},
		Body: template.HTML("<p>Body</p>"),
	}

	html, err := Render(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "collapsible_card.golden")
}
