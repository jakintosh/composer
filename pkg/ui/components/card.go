package components

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// CardProps describes a collapsible card with optional status + action elements and body.
type CardProps struct {
	Title  string
	Status *StatusBadgeProps
	Action *ButtonProps
	Body   g.Node
}

// Card renders the collapsible card structure used throughout the dashboard.
func Card(p CardProps) g.Node {
	summary := []g.Node{
		html.Span(
			html.Class("collapsible__title"),
			g.Text(p.Title),
		),
	}
	if p.Status != nil {
		summary = append(summary, StatusBadge(*p.Status))
	}
	if p.Action != nil {
		summary = append(summary, Button(*p.Action))
	}

	var body g.Node = g.Group(nil)
	if p.Body != nil {
		body = p.Body
	}

	return html.Li(
		html.Class("card card--collapsible"),
		html.Details(
			html.Class("collapsible"),
			html.Summary(
				html.Class("collapsible__summary"),
				g.Group(summary),
			),
			html.Div(
				html.Class("collapsible__content"),
				body,
			),
		),
	)
}
