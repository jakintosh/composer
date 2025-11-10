package columnheader

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"

	"composer/internal/ui/components/button"
)

// Props defines the content of a column header with optional action buttons.
type Props struct {
	Title   string
	Actions []button.Props
}

// Header renders the standard dashboard column header.
func Header(p Props) g.Node {
	actionNodes := make([]g.Node, 0, len(p.Actions))
	for _, action := range p.Actions {
		actionCopy := action
		actionNodes = append(actionNodes, button.Button(actionCopy))
	}

	return html.Header(
		html.Class("panel__header"),
		html.H2(html.Class("panel__title"), g.Text(p.Title)),
		html.Div(
			html.Class("panel__actions"),
			g.Group(actionNodes),
		),
	)
}
