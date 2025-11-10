package panel

import (
	"html/template"
	"testing"

	"composer/pkg/ui/components/columnheader"
	"gotest.tools/v3/golden"
)

func TestRenderSectionWithBody(t *testing.T) {
	props := SectionProps{
		Header:  columnheader.Props{Title: "Example"},
		Variant: VariantMuted,
		Body:    template.HTML("<p>Body</p>"),
	}

	html, err := RenderSection(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "section_with_body.golden")
}

func TestRenderSectionWithoutBody(t *testing.T) {
	props := SectionProps{
		Header: columnheader.Props{Title: "Empty"},
	}

	html, err := RenderSection(props)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	golden.Assert(t, string(html), "section_without_body.golden")
}
