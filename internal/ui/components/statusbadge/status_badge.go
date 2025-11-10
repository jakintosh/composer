package statusbadge

import (
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Props defines the label and modifier classes for a status badge.
type Props struct {
	Label   string
	Variant string
}

// Badge renders the status badge inline element.
func Badge(p Props) g.Node {
	className := "status-badge"
	if extra := strings.TrimSpace(p.Variant); extra != "" {
		className += " " + extra
	}
	return html.Span(
		html.Class(className),
		g.Text(p.Label),
	)
}
