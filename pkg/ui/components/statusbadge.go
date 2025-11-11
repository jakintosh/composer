package components

import (
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// StatusBadgeProps defines the label and modifier classes for a status badge.
type StatusBadgeProps struct {
	Label   string
	Variant string
}

// StatusBadge renders the status badge inline element.
func StatusBadge(p StatusBadgeProps) g.Node {
	className := "status-badge"
	if extra := strings.TrimSpace(p.Variant); extra != "" {
		className += " " + extra
	}
	return html.Span(
		html.Class(className),
		g.Text(p.Label),
	)
}
