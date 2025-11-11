package components

import (
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// ColumnVariant controls the background styling of the column surface.
type ColumnVariant string

const (
	ColumnVariantDefault ColumnVariant = ""
	ColumnVariantMuted   ColumnVariant = "panel--muted"
)

// ColumnProps configures a dashboard column with header metadata and list contents.
type ColumnProps struct {
	Title        string
	Actions      []ButtonProps
	Variant      ColumnVariant
	ListClass    string
	EmptyMessage string
	Items        []ColumnItem
}

// ColumnItem represents a single element rendered within the column list.
type ColumnItem struct {
	Class          string
	DisableWrapper bool
	Nodes          []g.Node
}

// ColumnSection renders the composed column with header + body list.
func ColumnSection(p ColumnProps) g.Node {
	className := "panel"
	if extra := strings.TrimSpace(string(p.Variant)); extra != "" {
		className += " " + extra
	}

	return html.Section(
		html.Class(className),
		ColumnHeader(p.Title, p.Actions),
		ColumnList(p.ListClass, p.EmptyMessage, p.Items),
	)
}

func ColumnHeader(title string, actions []ButtonProps) g.Node {
	actionNodes := make([]g.Node, 0, len(actions))
	for _, action := range actions {
		actionCopy := action
		actionNodes = append(actionNodes, Button(actionCopy))
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

func ColumnList(listClass, emptyMessage string, items []ColumnItem) g.Node {
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
		rendered = append(rendered, ColumnListItem(item))
	}

	return html.Ul(
		html.Class(className),
		g.Group(rendered),
	)
}

func ColumnListItem(item ColumnItem) g.Node {
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

// ColumnInfoRow renders a standard label/value pair used across column bodies.
func ColumnInfoRow(label, value string) g.Node {
	return html.P(
		html.Strong(g.Text(label+" ")),
		g.Text(value),
	)
}
