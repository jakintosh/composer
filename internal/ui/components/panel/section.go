package panel

import (
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"composer/internal/ui/components/columnheader"
)

// Variant controls the background styling of the panel container.
type Variant string

const (
	// VariantDefault renders the standard surface.
	VariantDefault Variant = ""
	// VariantMuted renders the muted surface used for the waiting column.
	VariantMuted Variant = "panel--muted"
)

// SectionProps describes a dashboard column with a header and optional body.
type SectionProps struct {
	Header  columnheader.Props
	Variant Variant
	Body    g.Node
}

// Section renders a panel container with the supplied props.
func Section(p SectionProps) g.Node {
	className := "panel"
	if extra := strings.TrimSpace(string(p.Variant)); extra != "" {
		className += " " + extra
	}

	children := []g.Node{
		columnheader.Header(p.Header),
	}
	if p.Body != nil {
		children = append(children, p.Body)
	}

	return html.Section(
		html.Class(className),
		g.Group(children),
	)
}
