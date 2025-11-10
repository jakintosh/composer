package datalist

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Item represents a row within the data list.
type Item struct {
	Primary   string
	Secondary g.Node
}

// Props describes a collection of rows rendered as a compact data list.
type Props struct {
	Items []Item
}

// List renders the configured rows. When no items are present the component
// renders nothing, matching the prior template behavior.
func List(p Props) g.Node {
	if len(p.Items) == 0 {
		return g.Group(nil)
	}

	rows := make([]g.Node, 0, len(p.Items))
	for _, item := range p.Items {
		var secondary g.Node = g.Group(nil)
		if item.Secondary != nil {
			secondary = item.Secondary
		}
		rows = append(rows, html.Li(
			html.Span(g.Text(item.Primary)),
			g.If(item.Secondary != nil, html.Span(secondary)),
		))
	}

	return html.Ul(
		html.Class("data-list"),
		g.Group(rows),
	)
}
