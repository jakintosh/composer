package panel

import (
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// ListProps defines the structure of a list within a panel.
type ListProps struct {
	Items        []ListItemProps
	ListClass    string
	EmptyMessage string
}

// ListItemProps describes a single row rendered inside a panel list.
type ListItemProps struct {
	Class          string
	Content        g.Node
	DisableWrapper bool
}

// List renders a panel list or an appropriate empty state.
func List(p ListProps) g.Node {
	if len(p.Items) == 0 {
		message := strings.TrimSpace(p.EmptyMessage)
		if message == "" {
			message = "No items."
		}
		return html.P(g.Text(message))
	}

	className := "panel__list"
	if extra := strings.TrimSpace(p.ListClass); extra != "" {
		className += " " + extra
	}

	itemNodes := make([]g.Node, 0, len(p.Items))
	for _, item := range p.Items {
		itemNodes = append(itemNodes, renderListItem(item))
	}

	return html.Ul(
		html.Class(className),
		g.Group(itemNodes),
	)
}

func renderListItem(item ListItemProps) g.Node {
	if item.DisableWrapper {
		if item.Content != nil {
			return item.Content
		}
		return g.Group(nil)
	}

	nodes := make([]g.Node, 0, 2)
	if className := strings.TrimSpace(item.Class); className != "" {
		nodes = append(nodes, html.Class(className))
	}
	if item.Content != nil {
		nodes = append(nodes, item.Content)
	}
	return html.Li(nodes...)
}
