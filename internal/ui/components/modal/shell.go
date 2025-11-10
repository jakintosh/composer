package modal

import (
	"strings"

	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

// Props defines the chrome rendered around a modal body.
type Props struct {
	ID         string
	Title      string
	TitleID    string
	CloseLabel string
	Body       g.Node
}

// DialogLabelID returns the DOM id used for aria-labelledby.
func (p Props) DialogLabelID() string {
	switch {
	case strings.TrimSpace(p.TitleID) != "":
		return strings.TrimSpace(p.TitleID)
	case strings.TrimSpace(p.ID) != "":
		return strings.TrimSpace(p.ID) + "-title"
	default:
		return "modal-title"
	}
}

// CloseLabelText returns the accessible label for the close button.
func (p Props) CloseLabelText() string {
	if label := strings.TrimSpace(p.CloseLabel); label != "" {
		return label
	}
	return "Close dialog"
}

// Shell renders the full modal chrome including the provided body.
func Shell(p Props) g.Node {
	var body g.Node = g.Group(nil)
	if p.Body != nil {
		body = p.Body
	}

	return html.Div(
		g.If(strings.TrimSpace(p.ID) != "", html.ID(strings.TrimSpace(p.ID))),
		html.Class("modal"),
		html.Aria("hidden", "true"),
		html.Div(
			html.Class("modal__dialog"),
			g.Attr("role", "dialog"),
			html.Aria("modal", "true"),
			html.Aria("labelledby", p.DialogLabelID()),
			html.Div(
				html.Class("modal__header"),
				html.H2(
					html.ID(p.DialogLabelID()),
					html.Class("modal__title"),
					g.Text(p.Title),
				),
				html.Button(
					html.Type("button"),
					html.Class("button button--ghost button--icon modal__close"),
					g.Attr("data-close-modal", ""),
					html.Aria("label", p.CloseLabelText()),
					g.Text("×"),
				),
			),
			html.Div(
				html.Class("modal__body"),
				body,
			),
		),
	)
}
