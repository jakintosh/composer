package modal

import (
	_ "embed"
	"html/template"

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

// Props contains the data required to render the run modal. It currently has no
// configurable fields but exists for future extensibility.
type Props struct{}

// Render executes the template for the run modal component.
func Render(p Props) (template.HTML, error) {
	return tmpl.Render(p)
}

// MustRender wraps Render and converts failures into HTML comments.
func MustRender(p Props) template.HTML {
	return templates.SafeHTML(Render(p))
}
