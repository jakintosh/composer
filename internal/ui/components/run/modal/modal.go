package modal

import (
	_ "embed"
	"html/template"

	modalshell "composer/internal/ui/components/ui/modal/shell"
	"composer/internal/ui/templates"
)

//go:embed modal.tmpl
var modalTemplate string

var tmpl = templates.New(
	"run_modal",
	"components/run/modal/modal.tmpl",
	modalTemplate,
	nil,
)

//go:embed body.tmpl
var bodyTemplate string

var bodyTmpl = templates.New(
	"run_modal_body",
	"components/run/modal/body.tmpl",
	bodyTemplate,
	nil,
)

// Props contains the data required to render the run modal. It currently has no
// configurable fields but exists for future extensibility.
type Props struct{}

// RenderShell renders the shared modal shell with the run form body.
func (p Props) RenderShell() template.HTML {
	return modalshell.MustRender(modalshell.Props{
		ID:         "run-modal",
		Title:      "Start Run",
		CloseLabel: "Close start run form",
		Body:       renderBody(),
	})
}

// Render executes the template for the run modal component.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender wraps Render and converts failures into HTML comments.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}

func renderBody() template.HTML {
	return templates.SafeHTML(bodyTmpl.Render(struct{}{}))
}
