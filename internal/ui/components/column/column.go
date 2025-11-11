package column

import (
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"composer/internal/ui/components/button"
)

// Variant controls the background styling of the column surface.
type Variant string

const (
	VariantDefault Variant = ""
	VariantMuted   Variant = "panel--muted"
)

// Props configures a dashboard column with header metadata and list contents.
type Props struct {
	Title        string
	Actions      []button.Props
	Variant      Variant
	ListClass    string
	EmptyMessage string
	Items        []Item
}

// Item represents a single element rendered within the column list.
type Item struct {
	Class          string
	DisableWrapper bool
	Nodes          []g.Node
}

// Section renders the composed column with header + body list.
func Section(p Props) g.Node {
	className := "panel"
	if extra := strings.TrimSpace(string(p.Variant)); extra != "" {
		className += " " + extra
	}

	return html.Section(
		html.Class(className),
		renderHeader(p.Title, p.Actions),
		renderList(p.ListClass, p.EmptyMessage, p.Items),
	)
}

func renderHeader(title string, actions []button.Props) g.Node {
	actionNodes := make([]g.Node, 0, len(actions))
	for _, action := range actions {
		actionCopy := action
		actionNodes = append(actionNodes, button.Button(actionCopy))
	}

	return html.Header(
		html.Class("panel__header"),
		html.H2(html.Class("panel__title"), g.Text(title)),
		html.Div(
			html.Class("panel__actions"),
			g.Group(actionNodes),
		),
	)
}

func renderList(listClass, emptyMessage string, items []Item) g.Node {
	if len(items) == 0 {
		message := strings.TrimSpace(emptyMessage)
		if message == "" {
			message = "No items."
		}
		return html.P(g.Text(message))
	}

	className := "panel__list"
	if extra := strings.TrimSpace(listClass); extra != "" {
		className += " " + extra
	}

	rendered := make([]g.Node, 0, len(items))
	for _, item := range items {
		rendered = append(rendered, renderListItem(item))
	}

	return html.Ul(
		html.Class(className),
		g.Group(rendered),
	)
}

func renderListItem(item Item) g.Node {
	if item.DisableWrapper {
		if item.Nodes != nil {
			return g.Group(item.Nodes)
		}
		return g.Group(nil)
	}

	attrs := []g.Node{}
	if className := strings.TrimSpace(item.Class); className != "" {
		attrs = append(attrs, html.Class(className))
	}
	attrs = append(attrs, g.Group(item.Nodes))
	return html.Li(attrs...)
}

// InfoRow renders a standard label/value pair used across column bodies.
func InfoRow(label, value string) g.Node {
	return html.P(
		html.Strong(g.Text(label+" ")),
		g.Text(value),
	)
}
